package sprint

import (
	"errors"
	"github.com/gorilla/schema"
	"net/http"
	"reflect"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"time"
)

var decoder = schema.NewDecoder()

func GetCreateSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	type RequestData struct {
		ProjectID uint `schema:"projectID,required"`
	}

	var requestData RequestData
	err := decoder.Decode(&requestData, r.URL.Query())
	if err != nil {
		return err
	}

	projectID := requestData.ProjectID

	c := createSprintDialog(projectID)
	return c.Render(r.Context(), w)
}

type sprintFormData struct {
	StartDate time.Time `schema:"start_date,required"`
	EndDate   time.Time `schema:"end_date,required"`
	Velocity  *float32  `schema:"velocity,required"`
	ProjectID uint      `schema:"project_id,required"`
}

// source: https://stackoverflow.com/questions/49285635/golang-gorilla-parse-date-with-specific-format-from-form
var timeConverter = func(value string) reflect.Value {
	layout := "2006-01-02"

	if v, err := time.Parse(layout, value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{}
}

func PostSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	return postSprint(w, r, database.GetDatabase())
}

func postSprint(w http.ResponseWriter, r *http.Request, db database.Database) error {
	var sprintFormData sprintFormData
	decoder.RegisterConverter(time.Time{}, timeConverter)
	err := decoder.Decode(&sprintFormData, r.PostForm)
	if err != nil {
		return err
	}

	sprints, err := db.GetSprintByProject(sprintFormData.ProjectID)
	if err != nil {
		return err
	}

	valid, validationErr := fieldValidation(sprintFormData, sprints)

	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, validationErr.Error(), http.StatusBadRequest)
		return nil
	}

	var sprint = &database.Sprint{
		Title:     "Sprint " + strconv.Itoa(len(sprints)),
		StartDate: sprintFormData.StartDate,
		EndDate:   sprintFormData.EndDate,
		Velocity:  sprintFormData.Velocity,
		ProjectID: sprintFormData.ProjectID,
	}

	err = db.CreateSprint(sprint)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", "/productbacklog?projectID="+strconv.Itoa(int(sprintFormData.ProjectID)))
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func fieldValidation(sprintToAdd sprintFormData, sprints []database.Sprint) (bool, error) {
	// validation: end date should be after start date
	if sprintToAdd.StartDate.After(sprintToAdd.EndDate) || sprintToAdd.StartDate.Equal(sprintToAdd.EndDate) {
		return false, errors.New("start date should be before end date")
	}

	// validation: start date should not be in the past or today
	if sprintToAdd.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return false, errors.New("start date should not be in the past")
	}

	// validation: sprint should not overlap with an existing one
	for _, s := range sprints {
		if (s.StartDate.Before(sprintToAdd.EndDate) || s.StartDate.Equal(sprintToAdd.EndDate)) &&
			(s.EndDate.After(sprintToAdd.StartDate) || s.EndDate.Equal(sprintToAdd.StartDate)) {
			return false, errors.New("sprint should not overlap with an existing one")
		}
	}

	return true, nil
}

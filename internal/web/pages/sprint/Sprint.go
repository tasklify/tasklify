package sprint

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"net/http"
	"reflect"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"
	"time"
)

var decoder = schema.NewDecoder()

type sprintFormData struct {
	StartDate time.Time `schema:"start_date,required"`
	EndDate   time.Time `schema:"end_date,required"`
	Velocity  *float32  `schema:"velocity,required"`
}

func GetCreateSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}

	c := createSprintDialog(uint(projectID))
	return c.Render(r.Context(), w)
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

	var sprintFormData sprintFormData
	decoder.RegisterConverter(time.Time{}, timeConverter)
	err := decoder.Decode(&sprintFormData, r.PostForm)
	if err != nil {
		return err
	}

	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))

	sprints, err := database.GetDatabase().GetSprintByProject(uint(projectID))
	if err != nil {
		return err
	}

	err2, fieldsInvalid := fieldValidation(w, r, sprintFormData, sprints)

	if fieldsInvalid {
		return err2
	}

	var sprint = &database.Sprint{
		Title:     "Sprint " + strconv.Itoa(len(sprints)),
		StartDate: sprintFormData.StartDate,
		EndDate:   sprintFormData.EndDate,
		Velocity:  sprintFormData.Velocity,
		ProjectID: uint(projectID),
	}

	if err != nil {
		return err
	}

	err = database.GetDatabase().CreateSprint(sprint)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", "/productbacklog?projectID="+strconv.Itoa(projectID))
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func fieldValidation(w http.ResponseWriter, r *http.Request, sprintFormData sprintFormData, sprints []database.Sprint) (error, bool) {

	// validation: end date before start date
	if sprintFormData.EndDate.Before(sprintFormData.StartDate) {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Start date should be before end date.")
		return c.Render(r.Context(), w), true
	}

	// validation: start date should not be in the past
	if sprintFormData.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Start date should not be in the past.")

		return c.Render(r.Context(), w), true
	}

	// validation: sprint should not overlap with an existing one
	for _, s := range sprints {
		// (StartA <= EndB) and (EndA >= StartB)
		if (s.StartDate.Before(sprintFormData.EndDate) || s.StartDate.Equal(sprintFormData.EndDate)) &&
			(s.EndDate.After(sprintFormData.StartDate) || s.EndDate.Equal(sprintFormData.StartDate)) {
			w.WriteHeader(http.StatusBadRequest)
			c := common.ValidationError("Sprint should not overlap with an existing one.")

			return c.Render(r.Context(), w), true
		}
	}

	// validation: sprint should not start or end on the weekend
	if sprintFormData.StartDate.Weekday() == time.Saturday ||
		sprintFormData.StartDate.Weekday() == time.Sunday ||
		sprintFormData.EndDate.Weekday() == time.Saturday ||
		sprintFormData.EndDate.Weekday() == time.Sunday {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Start/end date should not be on weekends.")

		return c.Render(r.Context(), w), true
	}
	return nil, false
}

func GetEditSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	sprintData, err := database.GetDatabase().GetSprintByID(uint(sprintID))
	if err != nil {
		return err
	}

	c := EditSprintDialog(*sprintData)
	return c.Render(r.Context(), w)
}

func PutSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	var sprintFormData sprintFormData
	decoder.RegisterConverter(time.Time{}, timeConverter)
	err := decoder.Decode(&sprintFormData, r.PostForm)
	if err != nil {
		return err
	}

	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))

	sprints, err := database.GetDatabase().GetSprintByProject(uint(projectID))
	if err != nil {
		return err
	}

	// remove current sprint from sprints
	n := 0
	for _, s := range sprints {
		if int(s.ID) != sprintID {
			sprints[n] = s
			n++
		}
	}

	sprints = sprints[:n]

	err2, fieldsInvalid := fieldValidation(w, r, sprintFormData, sprints)

	if fieldsInvalid {
		return err2
	}

	// get old sprint
	sprint, _ := database.GetDatabase().GetSprintByID(uint(sprintID))

	// change data
	sprint.EndDate = sprintFormData.EndDate
	sprint.StartDate = sprintFormData.StartDate
	sprint.Velocity = sprintFormData.Velocity

	// update sprint
	err = database.GetDatabase().UpdateSprint(sprint)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", "/productbacklog?projectID="+strconv.Itoa(projectID))
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func DeleteSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	// delete sprint
	err = database.GetDatabase().DeleteSprint(uint(projectID), uint(sprintID))

	w.Header().Set("HX-Redirect", "/productbacklog?projectID="+strconv.Itoa(projectID))
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

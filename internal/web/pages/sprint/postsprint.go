package sprint

import (
	"github.com/gorilla/schema"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"
	"time"
)

type sprintFormData struct {
	StartDate time.Time `schema:"start_date,required"`
	EndDate   time.Time `schema:"end_date,required"`
	Velocity  *float32  `schema:"velocity,required"`
}

var decoder = schema.NewDecoder()

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

	var projectID uint = 1 // TODO change

	err2, fieldsInvalid := fieldValidation(w, r, sprintFormData, projectID)

	if fieldsInvalid {
		return err2
	}

	var sprint = &database.Sprint{
		Title:     strconv.Itoa(rand.Int()), // TODO ask if title is needed for sprint
		StartDate: sprintFormData.StartDate,
		EndDate:   sprintFormData.EndDate,
		Velocity:  sprintFormData.Velocity,
		ProjectID: projectID, // Todo, when projects are implemented, change this
	}

	err = database.GetDatabase().CreateSprint(sprint)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", "/about")
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func fieldValidation(w http.ResponseWriter, r *http.Request, sprintFormData sprintFormData, projectID uint) (error, bool) {

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
	sprints, err := database.GetDatabase().GetSprintByProject(projectID)
	if err != nil {
		return err, false //TODO check if okay
	}

	for _, s := range sprints {
		// (StartA <= EndB) and (EndA >= StartB)
		if (s.StartDate.Before(sprintFormData.EndDate) || s.StartDate.Equal(sprintFormData.EndDate)) &&
			(s.EndDate.After(sprintFormData.StartDate) || s.EndDate.Equal(sprintFormData.StartDate)) {
			w.WriteHeader(http.StatusBadRequest)
			c := common.ValidationError("Sprint should not overlap with an existing one.")

			return c.Render(r.Context(), w), true
		}
	}

	return nil, false
}

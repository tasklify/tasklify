package userstory

import (
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"
	"tasklify/internal/web/pages"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.RegisterConverter(database.Priority{""}, priorityConverter)
}

func priorityConverter(value string) reflect.Value {
	priorityPtr := database.Priorities.Parse(value)
	if priorityPtr == nil {
		return reflect.Zero(reflect.TypeOf((*database.Priority)(nil)).Elem())
	}
	return reflect.ValueOf(*priorityPtr)
}

func PostUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type UserStoryFormData struct {
		Title           string            `schema:"title,required"`
		Description     string            `schema:"description,required"`
		Priority        database.Priority `schema:"priority,required"`
		BusinessValue   int               `schema:"business_value,required"`
		AcceptanceTests []string          `schema:"acceptanceTests"`
	}
	var userStoryData UserStoryFormData
	if err := decoder.Decode(&userStoryData, r.PostForm); err != nil {
		return err
	}

	ProjectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	// Check if a user story with the same title already exists
	userStoryExists := database.GetDatabase().UserStoryInThisProjectAlreadyExists(userStoryData.Title, uint(ProjectID))
	if userStoryExists {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("User story with same title already exists.")
		return c.Render(r.Context(), w)
	}
	if userStoryData.BusinessValue < 0 {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Business value must be a positive integer.")
		return c.Render(r.Context(), w)
	}

	userStory := &database.UserStory{
		Title:           userStoryData.Title,
		Description:     &userStoryData.Description,
		BusinessValue:   userStoryData.BusinessValue,
		Priority:        userStoryData.Priority,
		ProjectID:       uint(ProjectID),
		Realized:        new(bool), // Defaults to false
		AcceptanceTests: []database.AcceptanceTest{},
	}

	if err := database.GetDatabase().CreateUserStory(userStory); err != nil {
		return err
	}

	for _, testDescription := range userStoryData.AcceptanceTests {
		acceptanceTest := &database.AcceptanceTest{
			Description: &testDescription,
			Realized:    new(bool), // Defaults to false
			UserStoryID: userStory.ID,
		}
		if err := database.GetDatabase().CreateAcceptanceTest(acceptanceTest); err != nil {
			return err
		}
	}

	redirectURL := fmt.Sprintf("/productbacklog?projectID=%d", uint(ProjectID))
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func GetUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	ProjectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}

	projectRoles, err := database.GetDatabase().GetProjectRoles(params.UserID, uint(ProjectID))
	if err != nil {
		return err
	}

	if len(projectRoles) == 0 || slices.Contains(projectRoles, database.ProjectRoleDeveloper) {
		return pages.NotFound(w, r)
	}

	c := CreateUserStoryDialog(uint(ProjectID))
	return c.Render(r.Context(), w)
}

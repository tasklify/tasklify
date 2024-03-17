package userstory

import (
	"fmt"
	"net/http"
	"reflect"
	"tasklify/internal/database"
	"tasklify/internal/handlers"

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
		Title         string `schema:"title,required"`
		Description   string `schema:"description,required"`
		Priority      database.Priority  `schema:"priority,required"`
		BusinessValue int    `schema:"business_value,required"`
		ProjectID    uint   `schema:"projectID,required"`
		AcceptanceTests []string `schema:"acceptanceTests"`
	}
	var userStoryData UserStoryFormData
	if err := decoder.Decode(&userStoryData, r.PostForm); err != nil {
		return err
	}

	ProjectID := userStoryData.ProjectID

	userStory := &database.UserStory{
		Title:         userStoryData.Title,
		Description:   &userStoryData.Description,
		BusinessValue: userStoryData.BusinessValue,
		Priority:      userStoryData.Priority,
		ProjectID:     ProjectID,
		Realized:      new(bool), // Defaults to false
		AcceptanceTests: []database.AcceptanceTest{},
	}

	if err := database.GetDatabase().CreateUserStory(userStory); err != nil {
		return err
	}

	for _, testDescription := range userStoryData.AcceptanceTests {
		acceptanceTest := &database.AcceptanceTest{
			Description:   &testDescription,
			Realized:      new(bool), // Defaults to false
			UserStoryID:   userStory.ID,
		}
		if err := database.GetDatabase().CreateAcceptanceTest(acceptanceTest); err != nil {
			return err
		}
	}

    redirectURL := fmt.Sprintf("/productbacklog?projectID=%d", userStoryData.ProjectID)
    w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func GetUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		ProjectID uint `schema:"projectID,required"`
	}
	var requestData RequestData
	err := decoder.Decode(&requestData, r.URL.Query())
	if err != nil {
		return err
	}

	ProjectID := requestData.ProjectID

	c := CreateUserStoryDialog(ProjectID)
	return c.Render(r.Context(), w)
}

package userstory

import (
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

type UserStoryFormData struct {
	Title         string `schema:"title,required"`
	Description   string `schema:"description,required"`
	Priority      database.Priority  `schema:"priority,required"`
	BusinessValue int    `schema:"business_value,required"`
}

func PostUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	var userStoryData UserStoryFormData
	if err := decoder.Decode(&userStoryData, r.PostForm); err != nil {
		return err
	}

	projectID := uint(1)

	userStory := &database.UserStory{
		Title:         userStoryData.Title,
		Description:   &userStoryData.Description,
		BusinessValue: userStoryData.BusinessValue,
		Priority:      userStoryData.Priority,
		ProjectID:     projectID,
		Realized:      new(bool), // Defaults to false
	}

	if err := database.GetDatabase().CreateUserStory(userStory); err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", "/about")
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func GetUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	c := CreateUserStoryDialog()
	return c.Render(r.Context(), w)
}

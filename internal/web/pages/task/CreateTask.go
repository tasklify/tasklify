package task

import (
	"net/http"
	"tasklify/internal/database"
	"tasklify/internal/handlers"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func GetCreateTask(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	type RequestData struct {
		UserStoryID uint `schema:"userStoryID,required"`
	}

	var requestData RequestData
	err := decoder.Decode(&requestData, r.URL.Query())
	if err != nil {
		return err
	}

	UserStoryID := requestData.UserStoryID

	// TODO samo člani razvojne skupine
	users, err := database.GetDatabase().GetUsers()
	if err != nil {
		return err
	}

	c := createTaskDialog(UserStoryID, users)
	return c.Render(r.Context(), w)

}

type TaskFormData struct {
	Title        string   `schema:"title,required"`
	Description  string   `schema:"description,required"`
	TimeEstimate *float32 `schema:"time_estimate,required"`
	UserID       uint     `schema:"user_id"`
	UserStoryID  uint     `schema:"user_story_id,required"`
	ProjectID    uint     `schema:"project_id,required"`
}

func PostTask(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	// TODO

	return nil
}

package task

import (
	"fmt"
	"net/http"
	"tasklify/internal/database"
	"tasklify/internal/handlers"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func GetCreateTask(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	type RequestData struct {
		UserStoryID uint `schema:"userStoryID,required"`
		SprintID    uint `schema:"sprintID,required"`
		ProjectID   uint `schema:"projectID,required"`
	}

	var requestData RequestData
	err := decoder.Decode(&requestData, r.URL.Query())
	if err != nil {
		return err
	}

	userStoryID := requestData.UserStoryID
	sprintID := requestData.SprintID
	projectID := requestData.ProjectID

	// TODO Skrbnik metodologije in člani razvojne skupine
	users, err := database.GetDatabase().GetUsersOnProject(projectID)
	if err != nil {
		return err
	}

	c := createTaskDialog(projectID, sprintID, userStoryID, users)
	return c.Render(r.Context(), w)

}

type TaskFormData struct {
	Title        string   `schema:"title,required"`
	Description  string   `schema:"description,required"`
	TimeEstimate *float32 `schema:"time_estimate,required"`
	UserID       uint     `schema:"user_id"`
	UserStoryID  uint     `schema:"user_story_id,required"`
	ProjectID    uint     `schema:"project_id,required"`
	SprintID     uint     `schema:"sprint_id,required"`
}

func PostTask(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	var taskFormData TaskFormData
	if err := decoder.Decode(&taskFormData, r.PostForm); err != nil {
		return err
	}

	var task = &database.Task{
		Title:          &taskFormData.Title,
		Description:    &taskFormData.Description,
		TimeEstimate:   taskFormData.TimeEstimate,
		UserAccepted:   new(bool),
		Status:         &database.StatusTodo,
		ProjectID:      taskFormData.ProjectID,
		UserID:         nil,
		ProjectHasUser: nil,
		UserStoryID:    taskFormData.UserStoryID,
	}

	if taskFormData.UserID != 0 {
		projectHasUser, err := database.GetDatabase().GetProjectHasUserByProjectAndUser(taskFormData.UserID, taskFormData.ProjectID)
		if err != nil {
			return err
		}

		task.ProjectHasUser = projectHasUser
		task.UserID = &taskFormData.UserID
	}

	if err := database.GetDatabase().CreateTask(task); err != nil {
		return err
	}

	redirectURL := fmt.Sprintf("/sprintbacklog?sprintID=%d", taskFormData.SprintID)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

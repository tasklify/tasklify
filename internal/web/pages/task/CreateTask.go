package task

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func GetCreateTask(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	type RequestData struct {
		SprintID  uint `schema:"sprintID,required"`
		ProjectID uint `schema:"projectID,required"`
	}

	userStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}

	var requestData RequestData
	err1 := decoder.Decode(&requestData, r.URL.Query())
	if err1 != nil {
		return err
	}

	sprintID := requestData.SprintID
	projectID := requestData.ProjectID

	// only developers
	users, err := database.GetDatabase().GetUsersWithRoleOnProject(projectID, database.ProjectRoleDeveloper)
	if err != nil {
		return err
	}

	c := createTaskDialog(projectID, sprintID, uint(userStoryID), users)
	return c.Render(r.Context(), w)

}

type TaskFormData struct {
	Title        string   `schema:"title,required"`
	Description  string   `schema:"description"`
	TimeEstimate *float32 `schema:"time_estimate,required"`
	UserID       uint     `schema:"user_id"`
	ProjectID    uint     `schema:"project_id,required"`
	SprintID     uint     `schema:"sprint_id,required"`
}

func PostTask(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	var taskFormData TaskFormData
	if err := decoder.Decode(&taskFormData, r.PostForm); err != nil {
		return err
	}

	userStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
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
		UserStoryID:    uint(userStoryID),
	}

	if taskFormData.UserID != 0 {
		projectHasUser, err := database.GetDatabase().GetProjectHasUserByProjectAndUser(taskFormData.UserID, taskFormData.ProjectID)
		if err != nil {
			return err
		}

		task.ProjectHasUser = projectHasUser
		task.UserID = &taskFormData.UserID

		// if user created task and assigned to himself, the user should be automatically accepted
		if params.UserID == *task.UserID {
			userAccepted := true
			task.UserAccepted = &userAccepted
		}
	}

	if err := database.GetDatabase().CreateTask(task); err != nil {
		return err
	}

	redirectURL := fmt.Sprintf("/sprintbacklog/%d", taskFormData.SprintID)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

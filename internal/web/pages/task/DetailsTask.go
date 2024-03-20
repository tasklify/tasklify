package task

import (
	"net/http"
	"strconv"
	"tasklify/internal/database"

	"tasklify/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func GetTaskDetails(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	c := TaskDetailsDialog(*task)
	return c.Render(r.Context(), w)
}

func GetUserStoryFromTask(task database.Task) *database.UserStory {
	userStory, _ := database.GetDatabase().GetUserStoryByID(task.UserStoryID)
	return userStory
}
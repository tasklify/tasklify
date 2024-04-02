package sprintbacklog

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// GetSprintBacklog handles the request for fetching and displaying the sprint backlog.
func GetSprintBacklog(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	sprint, _ := database.GetDatabase().GetSprintByID(uint(sprintID))
	if sprint == nil {
		return pages.NotFound(w, r)
	}

	projectRoles, _ := database.GetDatabase().GetProjectRoles(params.UserID, sprint.ProjectID)
	if len(projectRoles) == 0 || slices.Contains(projectRoles, database.ProjectRoleManager) {
		return pages.NotFound(w, r)
	}

	sprintStatus := sprint.DetermineStatus()
	if sprintStatus != database.StatusInProgress {
		return pages.NotFound(w, r)
	}

	c := sprintBacklog(sprint, projectRoles, params.UserID)

	return pages.Layout(c, "Sprint Backlog", r).Render(r.Context(), w)
}

func GetUserFirstAndLastNameFromID(userID uint) string {
	user, _ := database.GetDatabase().GetUserByID(userID)
	return user.FirstName + " " + user.LastName
}

func mapTasksToStatuses(tasks []database.Task) (statusMap map[string][]database.Task) {
	statusMap = make(map[string][]database.Task)
	for _, task := range tasks {
		if task.UserID == nil {
			statusMap["Unassigned"] = append(statusMap["Unassigned"], task)
		} else {
			if !*task.UserAccepted {
				statusMap["Pending"] = append(statusMap["Pending"], task)
			} else if *task.Status == database.StatusInProgress {
				statusMap["Active"] = append(statusMap["Active"], task)
			} else if *task.Status == database.StatusDone {
				statusMap["Done"] = append(statusMap["Done"], task)
			} else {
				statusMap["Assigned"] = append(statusMap["Assigned"], task)
			}
		}
	}
	return
}

func GetSumOfTimeEstimates(tasks []database.Task) (sum float32) {
	sum = 0
	for _, task := range tasks {
		if task.TimeEstimate != nil {
			sum += *task.TimeEstimate
		}
	}
	return
}

func UnassignTask(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	task.UserAccepted = new(bool)
	task.UserID = nil

	err = database.GetDatabase().UpdateTask(task)

	w.Header().Set("HX-Redirect", fmt.Sprint("/sprintbacklog/", sprintID))
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

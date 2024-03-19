package sprintbacklog

import (
	"net/http"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// GetSprintBacklog handles the request for fetching and displaying the sprint backlog.
func GetSprintBacklog(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		SprintID uint `schema:"sprintID,required"`
	}
	var requestData RequestData
	err := decoder.Decode(&requestData, r.URL.Query())
	if err != nil {
		return err
	}

	sprintID := requestData.SprintID

	//fetch sprint
	sprint := database.GetDatabase().GetSprintByID(sprintID)
    if sprint == nil {
        return pages.NotFound(w, r)
    }

	c := sprintBacklog(sprint)

	return pages.Layout(c, "Sprint Backlog").Render(r.Context(), w)
}

func GetUsernameFromID(userID uint) string {
    user,_ := database.GetDatabase().GetUserByID(userID)
    return user.Username
}

func mapTasksToStatuses(tasks []database.Task) (statusMap map[string][]database.Task) {
	statusMap = make(map[string][]database.Task)
	for _, task := range tasks {
		if task.UserID == nil {
			statusMap["Unassigned"] = append(statusMap["Unassigned"], task)
		} else {
			if *task.Status == database.StatusInProgress {
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
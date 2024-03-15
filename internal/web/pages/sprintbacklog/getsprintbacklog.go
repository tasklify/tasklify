package sprintbacklog

import (
	"net/http"
	"sort"
	"tasklify/internal/database"
	"tasklify/internal/handlers"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// TaskWithUserStory enriches the Task with its User Story title.
type TaskWithUserStory struct {
    database.Task
    UserStoryTitle string
    AssignedTo string
}

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

    userStories, err := database.GetDatabase().GetUserStoriesBySprint(sprintID)
    if err != nil {
        return err
    }

    userStoryTitles, err := fetchUserStoryTitles(userStories)
    if err != nil {
        return err
    }

    allTasks, err := categorizeTasks(userStories, userStoryTitles)
    if err != nil {
        return err
    }



    // Get sort parameter from query
    sortParam := r.URL.Query().Get("sort")
    sortTasks(allTasks, sortParam)
	//get project id
	projectID := userStories[0].ProjectID


    return sprintBacklog(userStories, allTasks, projectID).Render(r.Context(), w)
}

// fetchUserStoryTitles creates a map of user story IDs to their titles.
func fetchUserStoryTitles(userStories []database.UserStory) (map[uint]string, error) {
    titles := make(map[uint]string)
    for _, us := range userStories {
        titles[us.ID] = us.Title
    }
    return titles, nil
}


// Assume this function categorizes tasks correctly and assigns the status string to each task.
func categorizeTasks(userStories []database.UserStory, titles map[uint]string) ([]TaskWithUserStory, error) {
    var allTasks []TaskWithUserStory
    for _, us := range userStories {
        tasks, err := database.GetDatabase().GetTasksByUserStory(us.ID)
        if err != nil {
            return nil, err
        }
        for _, task := range tasks {
            assignedTo := "Unassigned" // Default to Unassigned
            if task.UserID != nil {
                user, err := database.GetDatabase().GetUserByID(*task.UserID)
                if err == nil && user != nil {
                    // Safely assign username if user exists
                    assignedTo = user.Username
                }
            }
            allTasks = append(allTasks, TaskWithUserStory{
                Task:           task,
                UserStoryTitle: titles[task.UserStoryID],
                AssignedTo:     assignedTo,
            })
        }
    }
    return allTasks, nil
}

func sortTasks(allTasks []TaskWithUserStory, sortParam string) {
    switch sortParam {
    case "title":
        sort.Slice(allTasks, func(i, j int) bool {
            if allTasks[i].Title != nil && allTasks[j].Title != nil {
                return *allTasks[i].Title < *allTasks[j].Title
            }
            return false
        })
    case "status":
        sort.Slice(allTasks, func(i, j int) bool {
            return statusPriority(allTasks[i].Status) < statusPriority(allTasks[j].Status)
        })
    case "user_story":
        sort.Slice(allTasks, func(i, j int) bool {
            return allTasks[i].UserStoryTitle < allTasks[j].UserStoryTitle
        })
    case "assignee":
        sort.Slice(allTasks, func(i, j int) bool {
            return allTasks[i].AssignedTo < allTasks[j].AssignedTo
        })
    }

}


func statusPriority(status *database.Status) int {
    if status == nil || *status == (database.Status{}) {
        return -1 // or some value that represents "undefined"
    }
    switch *status {
    case database.StatusTodo:
        return 1
    case database.StatusInProgress:
        return 2
    case database.StatusDone:
        return 3
    default:
        return -1 // Handle unknown status
    }
}

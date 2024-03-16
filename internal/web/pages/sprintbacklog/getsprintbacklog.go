package sprintbacklog

import (
	"net/http"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"

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

    projectID := database.GetDatabase().GetSprintByID(sprintID).ProjectID


	c := sprintBacklog(userStories, allTasks, projectID)

    return pages.Layout(c, "Sprint Backlog").Render(r.Context(), w)
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

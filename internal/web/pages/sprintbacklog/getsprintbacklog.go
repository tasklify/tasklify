package sprintbacklog

import (
	"net/http"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"
)

type TaskWithUserStory struct {
    database.Task
    UserStoryTitle string
}

func GetSprintBacklog(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
    var sprintID uint = 1 // Example sprint ID
    userStories, err := database.GetDatabase().GetUserStoriesBySprint(sprintID)
    if err != nil {
        return err
    }

    // Map to store user story titles by ID for quick lookup
    userStoryTitles := make(map[uint]string)
    for _, us := range userStories {
        userStoryTitles[us.ID] = us.Title
    }

    taskGroups := map[string][]TaskWithUserStory{
        "unassigned": {},
        "assigned":   {},
        "completed":  {},
        "active":     {},
    }

    // Iterate over user stories to fetch tasks and assign user story titles
    for _, us := range userStories {
        tasks, err := database.GetDatabase().GetTasksByUserStory(us.ID)
        if err != nil {
            return err
        }

        for _, task := range tasks {
            taskWithUS := TaskWithUserStory{
                Task:           task,
                UserStoryTitle: userStoryTitles[task.UserStoryID],
            }

			status := *task.Status // Assuming task.Status is of type database.Status and is always non-nil

			switch status {
			case database.StatusTodo:
				if task.UserID == nil {
					taskGroups["unassigned"] = append(taskGroups["unassigned"], taskWithUS)
				} else {
					taskGroups["assigned"] = append(taskGroups["assigned"], taskWithUS)
				}
			case database.StatusInProgress:
				taskGroups["active"] = append(taskGroups["active"], taskWithUS)
			case database.StatusDone:
				taskGroups["completed"] = append(taskGroups["completed"], taskWithUS)
			}
		}
	}

	content := sprintBacklog(userStories, taskGroups["unassigned"], taskGroups["assigned"], taskGroups["completed"], taskGroups["active"])

	return pages.Layout(content, "Sprint Backlog").Render(r.Context(), w)
}

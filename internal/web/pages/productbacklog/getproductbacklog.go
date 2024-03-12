package productbacklog

import (
	"net/http"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"
)

func GetProductBacklog(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	// get all user stories
	var projectID uint = 1 // TODO change
	userStories, err := database.GetDatabase().GetUserStoriesByProject(projectID)
	if err != nil {
		return err
	}

	// unassigned, unrealized user stories
	var usInBacklog, usInSprint = filterBacklog(userStories)

	var userStoriesBySprint = groupUserStoriesBySprint(usInSprint)
	c := productBacklog(usInBacklog, userStoriesBySprint)

	return pages.Layout(c, "My website").Render(r.Context(), w)
}

func filterBacklog(userStories []database.UserStory) (inBacklog []database.UserStory, inSprint []database.UserStory) {

	for _, us := range userStories {
		if us.UserID == 0 && *us.Realized == false && *us.SprintID == 0 {
			inBacklog = append(inBacklog, us)
		} else {
			inSprint = append(inSprint, us)
		}
	}
	return
}

func groupUserStoriesBySprint(userStories []database.UserStory) (userStoriesBySprint map[string][]database.UserStory) {

	userStoriesBySprint = make(map[string][]database.UserStory)

	for _, us := range userStories {

		var sprintId = strconv.FormatUint(uint64(*us.SprintID), 10)

		if val, ok := userStoriesBySprint[sprintId]; ok {
			// add to existing slice
			var usSlice = val
			usSlice = append(usSlice, us)
			userStoriesBySprint[sprintId] = usSlice

		} else {
			// create a new slice
			sprintUserStories := []database.UserStory{us}
			userStoriesBySprint[sprintId] = sprintUserStories
		}
	}

	return
}

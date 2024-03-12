package productbacklog

import (
	"net/http"
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
	var unassignedUnrealizedUs = filterUnassignedUnrealized(userStories)

	c := productBacklog(unassignedUnrealizedUs)

	return pages.Layout(c, "My website").Render(r.Context(), w)
}

func filterUnassignedUnrealized(userStories []database.UserStory) (ret []database.UserStory) {

	for _, us := range userStories {
		if us.UserID == 0 && *us.Realized == false { // TODO check if works
			ret = append(ret, us)
		}
	}
	return
}

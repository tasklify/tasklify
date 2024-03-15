package productbacklog

import (
	"fmt"
	"net/http"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func GetProductBacklog(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		ProjectID uint `schema:"projectID,required"`
	}
	var requestData RequestData
	err := decoder.Decode(&requestData, r.URL.Query())
	if err != nil {
		return err
	}

	projectID := requestData.ProjectID
	// log the projectID
	fmt.Println(projectID)

	userStories, err := database.GetDatabase().GetUserStoriesByProject(projectID)
	if err != nil {
		return err
	}

	sprints, err := database.GetDatabase().GetSprintByProject(projectID)
	if err != nil {
		return err
	}

	// unassigned, unrealized user stories
	var usInBacklog, usInSprint = filterBacklog(userStories)

	var userStoriesBySprint = groupUserStoriesBySprint(usInSprint)
	var sprintMap = mapSprintsToSprintIds(sprints)
	var activityMap = mapActivityToSprints(sprints)

	c := productBacklog(usInBacklog, userStoriesBySprint, sprintMap, activityMap, projectID)

	return pages.Layout(c, "Backlog").Render(r.Context(), w)
}

func filterBacklog(userStories []database.UserStory) (inBacklog []database.UserStory, inSprint []database.UserStory) {

	for _, us := range userStories {
		if us.UserID == nil && *us.Realized == false && us.SprintID == nil {
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

func mapSprintsToSprintIds(sprints []database.Sprint) (sprintMap map[string]database.Sprint) {

	sprintMap = make(map[string]database.Sprint)


	for _, sprint := range sprints {
		var sprintID = strconv.FormatUint(uint64(sprint.ID), 10)

		sprintMap[sprintID] = sprint
	}

	return
}

func mapActivityToSprints(sprints []database.Sprint) (activityMap map[string]database.Status) {

	activityMap = make(map[string]database.Status)


	for _, sprint := range sprints {
		var sprintID = strconv.FormatUint(uint64(sprint.ID), 10)

		activityMap[sprintID], _ = sprint.DetermineStatus()
	}

	return
}
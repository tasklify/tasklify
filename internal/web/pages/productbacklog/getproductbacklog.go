package productbacklog

import (
	"errors"
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
	var usInBacklog, _ = filterBacklog(userStories)
	var sprintMap = mapSprintsToSprintIds(sprints)

	c := productBacklog(usInBacklog, sprintMap, projectID)
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

func mapSprintsToSprintIds(sprints []database.Sprint) (sprintMap map[string]database.Sprint) {

	sprintMap = make(map[string]database.Sprint)

	for _, sprint := range sprints {
		var sprintID = strconv.FormatUint(uint64(sprint.ID), 10)

		sprintMap[sprintID] = sprint
	}

	return
}

func PostAddUserStoryToSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return err
	}

	usIDCount := len(r.Form["selectedTasks"])
	selectedTaskIds := make([]uint, 0, usIDCount)

	for _, id := range r.Form["selectedTasks"] {
		if usID, err := strconv.Atoi(id); err == nil {
			selectedTaskIds = append(selectedTaskIds, uint(usID))
		}
	}

	sprintID, err := strconv.Atoi(r.FormValue("sprintID"))
	if err != nil {
		return err
	}

	_, err = database.GetDatabase().AddUserStoryToSprint(uint(sprintID), selectedTaskIds)
	if err != nil {
		return err
	}

	fmt.Println("Sprint ID:", sprintID)
	fmt.Println("User Story IDs:", selectedTaskIds)

	callbackURL := r.FormValue("callback")
	if callbackURL != "" {
		w.Header().Set("HX-Redirect", callbackURL)
	} else {
		return errors.New("callback URL not provided")
	}

	w.WriteHeader(http.StatusSeeOther)
	return nil
}

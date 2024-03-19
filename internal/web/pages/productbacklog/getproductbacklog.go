package productbacklog

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
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

	userStories, err := database.GetDatabase().GetUserStoriesByProject(projectID)
	if err != nil {
		return err
	}

	sprints, err := database.GetDatabase().GetSprintByProject(projectID)

	// order sprints by start day
	sort.Slice(sprints, func(i, j int) bool {
		return sprints[i].StartDate.Before(sprints[j].StartDate)
	})

	if err != nil {
		return err
	}

	// unassigned, unrealized user stories
	var usInBacklog, _ = filterBacklog(userStories)

	//get user project role
	projectRole := database.GetDatabase().GetProjectRole(params.UserID, projectID)

	c := productBacklog(usInBacklog, sprints, projectID, projectRole)
	return pages.Layout(c, "Backlog", r).Render(r.Context(), w)
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

func RemoveUserStoryFromSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	if err := r.ParseForm(); err != nil {
		return err
	}

	userStoryID, err := strconv.Atoi(r.FormValue("userStoryID"))
	if err != nil {
		return err
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(uint(userStoryID))
	if err != nil {
		return err
	}

	// remove sprint from userStory
	userStory.SprintID = nil
	if err := database.GetDatabase().UpdateUserStory(userStory); err != nil {
		fmt.Println("Error updating user story")
		return err
	}

	callbackURL := r.FormValue("callback")
	if callbackURL != "" {
		w.Header().Set("HX-Redirect", callbackURL)
	} else {
		return errors.New("callback URL not provided")
	}

	w.WriteHeader(http.StatusSeeOther)
	return nil
}

func PostUserStoryAccepted(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	userStoryID, err := strconv.Atoi(r.FormValue("userStoryID"))
	if err != nil {
		return err
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(uint(userStoryID))
	if err != nil {
		return err
	}

	*userStory.Realized = true
	if err := database.GetDatabase().UpdateUserStory(userStory); err != nil {
		fmt.Println("Error updating user story")
		return err
	}

	callbackURL := r.FormValue("callback")
	fmt.Println(callbackURL)
	if callbackURL != "" {
		w.Header().Set("HX-Redirect", callbackURL)
	} else {
		return errors.New("callback URL not provided")
	}

	w.WriteHeader(http.StatusSeeOther)
	return nil
}

func PostUserStoryRejected(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	userStoryID, err := strconv.Atoi(r.FormValue("userStoryID"))
	if err != nil {
		return err
	}

	projectID, err := strconv.Atoi(r.FormValue("projectID"))
	if err != nil {
		return err
	}

	c := CreateRejectionCommentDialog(uint(userStoryID), "/productbacklog?projectID="+strconv.Itoa(int(projectID)))

	return c.Render(r.Context(), w)

}

func PostRejectionComment(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	userStoryID, err := strconv.Atoi(r.FormValue("userStoryID"))
	if err != nil {
		return err
	}

	comment := r.FormValue("comment")

	userStory, err := database.GetDatabase().GetUserStoryByID(uint(userStoryID))
	if err != nil {
		return err
	}

	userStory.RejectionComment = &comment
	userStory.SprintID = nil
	*userStory.Realized = false
	if err := database.GetDatabase().UpdateUserStory(userStory); err != nil {
		fmt.Println("Error updating user story")
		return err
	}

	callbackURL := r.FormValue("callback")
	fmt.Println(callbackURL)
	if callbackURL != "" {
		w.Header().Set("HX-Redirect", callbackURL)
	} else {
		return errors.New("callback URL not provided")
	}

	w.WriteHeader(http.StatusSeeOther)
	return nil
}

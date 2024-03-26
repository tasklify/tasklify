package productbacklog

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"
	"tasklify/internal/web/pages"

	"github.com/go-chi/chi/v5"
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

	project, err := database.GetDatabase().GetProjectByID(projectID)
	if err != nil {
		return err
	}

	userStories, err := database.GetDatabase().GetUserStoriesByProject(projectID)
	if err != nil {
		return err
	}

	sprints, err := database.GetDatabase().GetSprintByProject(projectID)

	// order sprints by create date - active should always be on top
	sort.Slice(sprints, func(i, j int) bool {
		statusI := sprints[i].DetermineStatus()
		statusJ := sprints[j].DetermineStatus()

		if statusI == database.StatusInProgress {
			return true
		} else if statusJ == database.StatusInProgress {
			return false
		}

		return sprints[i].CreatedAt.After(sprints[j].StartDate)
	})

	if err != nil {
		return err
	}

	// unassigned, unrealized user stories
	var usInBacklog, usInFuture, usInDone = filterBacklog(userStories)

	//get user project role
	projectRoles, _ := database.GetDatabase().GetProjectRoles(params.UserID, projectID)

	user, err := database.GetDatabase().GetUserByID(params.UserID)
	if err != nil {
		return err
	}

	c := productBacklog(usInBacklog, usInDone, usInFuture, sprints, projectID, projectRoles, *project, user.SystemRole)
	return pages.Layout(c, "Backlog", r).Render(r.Context(), w)
}

func filterBacklog(userStories []database.UserStory) (inBacklog, inFuture, inDone []database.UserStory) {
	for _, us := range userStories {
		if *us.Realized {
			inDone = append(inDone, us)
		} else if us.SprintID == nil {
			if us.Priority == database.PriorityWontHaveThisTime {
				inFuture = append(inFuture, us)
			} else if us.UserID == nil && *us.Realized == false && us.SprintID == nil {
				inBacklog = append(inBacklog, us)
			}
		}
	}
	return inBacklog, inFuture, inDone
}

func PostAddUserStoryToSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	// Parse form data
	type UserStoriesToSprintFormData struct {
		SprintID         uint     `schema:"sprintID,required"`
		Velocity         *float32 `schema:"velocity,required"`
		ActiveSprintLoad *float32 `schema:"active_sprint_load,required"`
		Callback         string   `schema:"callback,required"`
		UserStoryIDs     []uint   `schema:"selectedTasks"`
	}

	var userStoriesToSprintFormData UserStoriesToSprintFormData
	if err := decoder.Decode(&userStoriesToSprintFormData, r.PostForm); err != nil {
		return err
	}

	sprintID := userStoriesToSprintFormData.SprintID
	selectedUserStoryIds := userStoriesToSprintFormData.UserStoryIDs
	velocity := userStoriesToSprintFormData.Velocity
	activeSprintLoad := userStoriesToSprintFormData.ActiveSprintLoad
	currentLoad, err := database.GetDatabase().GetUserStoriesLoad(selectedUserStoryIds)

	if float32(currentLoad) > (*velocity - *activeSprintLoad) {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Sum of user story points exceed sprint velocity.")
		return c.Render(r.Context(), w)
	}

	_, err = database.GetDatabase().AddUserStoryToSprint(sprintID, selectedUserStoryIds)
	if err != nil {
		return err
	}

	fmt.Println("Sprint ID:", sprintID)
	fmt.Println("User Story IDs:", selectedUserStoryIds)

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

	userStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
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
	userStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
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

	w.Header().Set("HX-Redirect", "/productbacklog?projectID="+strconv.Itoa(int(userStory.ProjectID)))

	w.WriteHeader(http.StatusSeeOther)
	return nil
}

func GetUserStoryRejected(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	userStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}

	userStory, _ := database.GetDatabase().GetUserStoryByID(uint(userStoryID))

	projectRoles, err := database.GetDatabase().GetProjectRoles(params.UserID, uint(userStory.ProjectID))
	if err != nil {
		return err
	}

	if !slices.Contains(projectRoles, database.ProjectRoleManager) {
		return pages.NotFound(w, r)
	}

	sprint, _ := database.GetDatabase().GetSprintByID(*userStory.SprintID)
	if (*userStory.Realized) || (sprint.DetermineStatus() != database.StatusDone) {
		return pages.NotFound(w, r)
	}

	c := CreateRejectionCommentDialog(uint(userStoryID))

	return c.Render(r.Context(), w)

}

func PostUserStoryRejected(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RejectFormData struct {
		Comment string `schema:"comment,required"`
	}
	var rejectFormData RejectFormData
	if err := decoder.Decode(&rejectFormData, r.PostForm); err != nil {
		return err
	}
	userStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}

	comment := rejectFormData.Comment

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

	w.Header().Set("HX-Redirect", "/productbacklog?projectID="+strconv.Itoa(int(userStory.ProjectID)))

	w.WriteHeader(http.StatusSeeOther)
	return nil
}

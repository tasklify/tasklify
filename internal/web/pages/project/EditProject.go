package project

import (
	"fmt"
	"net/http"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"
	"tasklify/internal/web/pages"

	"github.com/go-chi/chi/v5"
)

func GetEditProjectInfo(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	project, err := database.GetDatabase().GetProjectByID(projectID)
	if err != nil {
		return err
	}

	c := editProjectInfo(*project)
	return pages.Layout(c, "Edit project", r).Render(r.Context(), w)
}

func UpdateProjectInfo(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		ProjectID   uint   `schema:"projectID,required"`
		Title       string `schema:"title,required"`
		Description string `schema:"description"`
	}

	var req RequestData
	err := decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	// Check if a project with the same title already exists
	projectExists := database.GetDatabase().ProjectWithTitleExists(req.Title, &req.ProjectID)
	if projectExists {
		fmt.Println("Exists")
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Project with the same title already exists.")
		return c.Render(r.Context(), w)
	}

	updatedProject := database.Project{
		Title:       req.Title,
		Description: req.Description,
	}

	err = database.GetDatabase().UpdateProject(req.ProjectID, updatedProject)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", fmt.Sprint("/project-info/", req.ProjectID))
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

var projectDevelopers = make(map[uint]database.User)
var productOwner database.User
var scrumMaster database.User

func GetEditProjectMembers(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	projectDevelopersSlice, err := database.GetDatabase().GetUsersWithRoleOnProject(projectID, database.ProjectRoleDeveloper)
	if err != nil {
		return err
	}

	projectDevelopers = make(map[uint]database.User)

	for _, developer := range projectDevelopersSlice {
		projectDevelopers[developer.ID] = developer
	}

	productOwners, err := database.GetDatabase().GetUsersWithRoleOnProject(projectID, database.ProjectRoleManager)
	if err != nil {
		return err
	}
	productOwner = productOwners[0]

	scrumMasters, err := database.GetDatabase().GetUsersWithRoleOnProject(projectID, database.ProjectRoleMaster)
	if err != nil {
		return err
	}
	scrumMaster = scrumMasters[0]

	availableUsersList, err := getFilteredUserList()
	if err != nil {
		return err
	}

	allUsersList, err := database.GetDatabase().GetUsers()
	if err != nil {
		return err
	}

	c := editProjectMembers(projectID, &productOwner, &scrumMaster, projectDevelopers, availableUsersList, allUsersList)
	return pages.Layout(c, "Edit project", r).Render(r.Context(), w)
}

func EditProjectRemoveDeveloper(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		UserID uint `schema:"userID,required"`
	}

	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	var req RequestData
	err = decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	delete(projectDevelopers, req.UserID)

	userList, err := getFilteredUserList()
	if err != nil {
		return err
	}

	c := ProjectDevelopersContainer(projectID, projectDevelopers, userList)
	return c.Render(r.Context(), w)
}

func EditProjectAddDeveloper(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		UserID uint `schema:"userID,required"`
	}

	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	var req RequestData
	err = decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	user, err := database.GetDatabase().GetUserByID(req.UserID)
	if err != nil {
		return err
	}

	projectDevelopers[req.UserID] = *user

	userList, err := getFilteredUserList()
	if err != nil {
		return err
	}

	c := ProjectDevelopersContainer(projectID, projectDevelopers, userList)
	return c.Render(r.Context(), w)

}

func EditProjectChangeOwner(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		UserID uint `schema:"productOwner,required"`
	}

	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	var req RequestData
	err = decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	user, err := database.GetDatabase().GetUserByID(req.UserID)
	if err != nil {
		return err
	}

	productOwner = *user

	allUsersList, err := database.GetDatabase().GetUsers()
	if err != nil {
		return err
	}

	fmt.Println(productOwner)

	c := ProductOwnerSelect(projectID, &productOwner, allUsersList)
	return c.Render(r.Context(), w)
}

func EditProjectChangeMaster(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		UserID uint `schema:"scrumMaster,required"`
	}

	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	var req RequestData
	err = decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	user, err := database.GetDatabase().GetUserByID(req.UserID)
	if err != nil {
		return err
	}

	scrumMaster = *user

	allUsersList, err := database.GetDatabase().GetUsers()
	if err != nil {
		return err
	}

	c := ScrumMasterSelect(projectID, &scrumMaster, allUsersList)
	return c.Render(r.Context(), w)
}

func UpdateProjectMembers(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		ProjectID uint `schema:"projectID,required"`
	}

	var req RequestData
	err := decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	if scrumMaster.ID == productOwner.ID {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Product owner and SCRUM master should not be the same person!")
		return c.Render(r.Context(), w)
	}

	if _, ok := projectDevelopers[productOwner.ID]; ok {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(fmt.Sprintf("User %v cannot be Product owner and Project developer at the same time.", productOwner.FirstName+" "+productOwner.LastName))
		return c.Render(r.Context(), w)
	}

	// if _, ok := projectDevelopers[scrumMaster.ID]; ok {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	c := common.ValidationError(fmt.Sprintf("User %v cannot be SCRUM master and Project developer at the same time.", productOwner.FirstName+" "+productOwner.LastName))
	// 	return c.Render(r.Context(), w)
	// }

	var addedUserIDs []uint

	projectData := database.Project{
		ProductOwnerID: productOwner.ID,
		ScrumMasterID:  scrumMaster.ID,
	}

	// Add project owner and scrum master
	if err := database.GetDatabase().UpdateProject(req.ProjectID, projectData); err != nil {
		return err
	}

	// Add all project developers
	for userID := range projectDevelopers {
		if err := database.GetDatabase().AddDeveloperToProject(req.ProjectID, userID); err != nil {
			return err
		}
		addedUserIDs = append(addedUserIDs, userID)
	}

	fmt.Println(addedUserIDs)

	// Remove all other users from project
	if err := database.GetDatabase().RemoveUsersNotInList(req.ProjectID, addedUserIDs); err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", fmt.Sprint("/project-info/", req.ProjectID))
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func getFilteredUserList() ([]database.User, error) {
	var userIDs []uint
	for userID := range projectDevelopers {
		userIDs = append(userIDs, userID)
	}

	if len(userIDs) > 0 {
		userList, err := database.GetDatabase().GetFilteredUsers(userIDs)
		if err != nil {
			return []database.User{}, err
		}

		return userList, nil
	} else {
		userList, err := database.GetDatabase().GetUsers()
		if err != nil {
			return []database.User{}, err
		}

		return userList, nil
	}
}

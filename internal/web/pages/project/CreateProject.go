package project

import (
	"fmt"
	"net/http"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func GetCreateProject(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	users, err := database.GetDatabase().GetUsers(nil)
	if err != nil {
		return err
	}
	c := createProjectDialog(users)
	return c.Render(r.Context(), w)
}

func PostCreateProject(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		Title          string `schema:"title,required"`
		Description    string `schema:"description"`
		ProductOwnerID uint   `schema:"productOwner,required"`
		ScrumMasterID  uint   `schema:"scrumMaster,required"`
	}

	var req RequestData
	err := decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	// Check if a project with the same title already exists
	projectExists := database.GetDatabase().ProjectWithTitleExists(req.Title, nil)
	if projectExists {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Project with the same title already exists.")
		return c.Render(r.Context(), w)
	}

	// Product owner should not be SCRUM master
	if req.ProductOwnerID == req.ScrumMasterID {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Product owner and SCRUM master should not be the same person!")
		return c.Render(r.Context(), w)
	}

	newProject := database.Project{
		Title:       req.Title,
		Description: req.Description,
	}

	pID, err := database.GetDatabase().CreateProject(&newProject)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Add Product owner
	if err := database.GetDatabase().AddUserToProject(pID, req.ProductOwnerID, database.ProjectRoleManager.Val); err != nil {
		return err
	}

	// Add SCRUM master
	if err := database.GetDatabase().AddUserToProject(pID, req.ScrumMasterID, database.ProjectRoleMaster.Val); err != nil {
		return err
	}

	users, err := database.GetDatabase().GetUsersNotOnProject(pID)
	if err != nil {
		return err
	}

	c := addProjectDevelopersDialog(pID, users, []database.User{})
	return c.Render(r.Context(), w)
}

func PostAddProjectDeveloper(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		ProjectID uint `schema:"projectID,required"`
		UserID    uint `schema:"userID,required"`
	}

	var req RequestData
	err := decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	if err := database.GetDatabase().AddUserToProject(req.ProjectID, req.UserID, database.ProjectRoleDeveloper.Val); err != nil {
		return err
	}

	users, err := database.GetDatabase().GetUsersNotOnProject(req.ProjectID)
	if err != nil {
		return err
	}

	projectDevelopers, err := database.GetDatabase().GetUsersWithRoleOnProject(req.ProjectID, database.ProjectRoleDeveloper)
	if err != nil {
		return err
	}

	c := addProjectDevelopersDialog(req.ProjectID, users, projectDevelopers)
	return c.Render(r.Context(), w)
}

func RemoveProjectDeveloper(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		ProjectID uint `schema:"projectID,required"`
		UserID    uint `schema:"userID,required"`
	}

	var requestData RequestData
	err := decoder.Decode(&requestData, r.PostForm)
	if err != nil {
		return err
	}

	if err := database.GetDatabase().RemoveUserFromProject(requestData.ProjectID, requestData.UserID); err != nil {
		return err
	}

	users, err := database.GetDatabase().GetUsersNotOnProject(requestData.ProjectID)
	if err != nil {
		return err
	}

	projectDevelopers, err := database.GetDatabase().GetUsersWithRoleOnProject(requestData.ProjectID, database.ProjectRoleDeveloper)
	if err != nil {
		return err
	}

	c := addProjectDevelopersDialog(requestData.ProjectID, users, projectDevelopers)
	return c.Render(r.Context(), w)
}

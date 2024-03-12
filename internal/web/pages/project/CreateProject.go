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

type ProjectUserDetails struct {
	Username    string
	Email       string
	FirstName   string
	LastName    string
	ProjectRole string
}

func GetCreateProject(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	// users, err := database.GetDatabase().GetUsers()
	// if err != nil {
	// 	return err
	// }

	// signedInUser, err := database.GetDatabase().GetUserByID(params.UserID)
	// if err != nil {
	// 	return err
	// }

	// d := Data{
	// 	signedInUserID: signedInUser.ID,
	// 	users:          users,
	// 	projectMembers: []database.User{*signedInUser, *signedInUser},
	// }

	c := createProjectDialog()
	return c.Render(r.Context(), w)
}

func PostCreateProject(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	var projectData database.Project
	err := decoder.Decode(&projectData, r.PostForm)
	if err != nil {
		return err
	}

	// Check if a project with the same title already exists
	projectExists := database.GetDatabase().ProjectWithTitleExists(projectData.Title)
	if projectExists {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Project with same title already exists.")
		return c.Render(r.Context(), w)
	}

	pID, err := database.GetDatabase().CreateProject(&projectData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	users, err := database.GetDatabase().GetUsers()
	if err != nil {
		return err
	}

	c := addProjectMembersDialog(pID, users, []database.User{})
	return c.Render(r.Context(), w)
}

func PostAddProjectMember(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type ProjectMemberData struct {
		ProjectID uint   `schema:"projectID,required"`
		UserID    uint   `schema:"userID,required"`
		RoleID    string `schema:"roleID,required"`
	}

	var projectMemberData ProjectMemberData
	err := decoder.Decode(&projectMemberData, r.PostForm)
	if err != nil {
		return err
	}

	if err := database.GetDatabase().AddUserToProject(projectMemberData.ProjectID, projectMemberData.UserID, projectMemberData.RoleID); err != nil {
		return err
	}

	users, err := database.GetDatabase().GetUsersNotOnProject(projectMemberData.ProjectID)
	if err != nil {
		return err
	}

	projectMembers, err := database.GetDatabase().GetUsersOnProject(projectMemberData.ProjectID)
	if err != nil {
		return err
	}

	c := addProjectMembersDialog(projectMemberData.ProjectID, users, projectMembers)
	return c.Render(r.Context(), w)
}

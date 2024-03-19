package project

import (
	"fmt"
	"net/http"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"
	"tasklify/internal/web/pages"

	"github.com/go-chi/chi"
)

// var projectDevelopers []database.User
// var productOwner database.User
// var scrumMaster database.User

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

	// projectDevelopers, err = database.GetDatabase().GetUsersWithRoleOnProject(projectID, database.ProjectRoleDeveloper)
	// if err != nil {
	// 	return err
	// }

	// productOwners, err := database.GetDatabase().GetUsersWithRoleOnProject(projectID, database.ProjectRoleManager)
	// if err != nil {
	// 	return err
	// }

	// if len(productOwners) > 0 {
	// 	productOwner = productOwners[0]
	// }

	// scrumMasters, err := database.GetDatabase().GetUsersWithRoleOnProject(projectID, database.ProjectRoleMaster)
	// if err != nil {
	// 	return err
	// }

	// if len(scrumMasters) > 0 {
	// 	scrumMaster = scrumMasters[0]
	// }

	c := editProjectInfo(*project)
	return pages.Layout(c, "Edit project").Render(r.Context(), w)
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

	w.Header().Set("HX-Redirect", "/productbacklog?projectID="+strconv.Itoa(int(req.ProjectID)))
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

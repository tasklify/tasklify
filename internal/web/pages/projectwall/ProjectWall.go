package projectwall

import (
	"fmt"
	"net/http"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func GetProjectWall(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	project, err := database.GetDatabase().GetProjectByID(projectID)
	if err != nil {
		return err
	}

	projectRoles, _ := database.GetDatabase().GetProjectRoles(params.UserID, projectID)
	if len(projectRoles) == 0 {
		return pages.NotFound(w, r)
	}

	user, err := database.GetDatabase().GetUserByID(params.UserID)
	if err != nil {
		return err
	}

	posts, err := database.GetDatabase().GetProjectWallPosts(projectID)
	if err != nil {
		return err
	}
	project.ProjectWallPosts = posts

	c := projectWall(*project, projectRoles, *user)
	return pages.Layout(c, "Project Wall", r).Render(r.Context(), w)
}

func GetNewPostDialog(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	c := AddNewPostDialog(uint(projectID))
	return c.Render(r.Context(), w)
}

func AddNewPost(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		Body string `schema:"body,required"`
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

	newPost := database.ProjectWallPost{
		ProjectID: projectID,
		AuthorID:  params.UserID,
		Body:      req.Body,
	}

	err = database.GetDatabase().AddProjectWallPost(newPost)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", fmt.Sprint("/project-wall/", projectID))

	return nil
}

func GetEditPost(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	postIDInt, err := strconv.Atoi(chi.URLParam(r, "postID"))
	if err != nil {
		return err
	}
	postID := uint(postIDInt)

	postData, err := database.GetDatabase().GetProjectWallPostByID(postID)
	if err != nil {
		return err
	}

	c := EditPostDialog(*postData)
	return c.Render(r.Context(), w)
}

func PutPost(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		Body string `schema:"body,required"`
	}

	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	postIDInt, err := strconv.Atoi(chi.URLParam(r, "postID"))
	if err != nil {
		return err
	}
	postID := uint(postIDInt)

	var req RequestData
	err = decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	err = database.GetDatabase().EditProjectWallPost(projectID, postID, req.Body)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", fmt.Sprint("/project-wall/", projectID))

	return nil
}

func DeletePost(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	projectIDInt, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}
	projectID := uint(projectIDInt)

	postIDInt, err := strconv.Atoi(chi.URLParam(r, "postID"))
	if err != nil {
		return err
	}
	postID := uint(postIDInt)

	err = database.GetDatabase().DeleteProjectWallPost(projectID, postID)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", fmt.Sprint("/project-wall/", projectID))

	return nil
}

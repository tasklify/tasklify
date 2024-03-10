package project

import (
	"fmt"
	"net/http"
	"tasklify/internal/handlers"
)

type PostProjectHandler struct{}

func NewPostProjectHandler() *PostProjectHandler {
	return &PostProjectHandler{}
}

type ProjectRequest struct {
	Title       string
	Users       []ProjectUserDetails
	Sprints     []uint
	UserStories []uint
}

type ProjectUserDetails struct {
	Username    string
	ProjectRole string
}

func PostProject(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	fmt.Println(r.FormValue("projectUser"))
	// reqData := &ProjectRequest{}
	// json.NewDecoder(r.Body).Decode(&reqData)

	// newProject := database.Project{
	// 	Title: reqData.Title,
	// 	Users: []database.User{},
	// }

	// fmt.Println(reqData)

	// for _, userDetail := range reqData.Users {
	// 	// user, _ := database.GetDatabase().GetUser(userDetail.Username)
	// 	// newProject.Users = append(newProject.Users, *user)
	// 	fmt.Println(userDetail)
	// }
	// fmt.Println(newProject)
	// // database.GetDatabase().CreateProject(&newProject)

	return nil
}

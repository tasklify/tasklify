package project

import (
	"net/http"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"
)

type GetCreateProjectHandler struct{}

func NewGetCreateProjectHandler() *GetCreateProjectHandler {
	return &GetCreateProjectHandler{}
}

type Data struct {
	users []database.User
}

func GetCreateProject(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	users, err := database.GetDatabase().GetUsers()

	if err != nil {
		return err
	}

	d := Data{
		users: users,
	}

	c := createProjectDialog(d)
	return pages.Layout(c, "Tasklify").Render(r.Context(), w)
}

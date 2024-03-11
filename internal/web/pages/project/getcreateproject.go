package project

import (
	"net/http"
	"tasklify/internal/handlers"
)

func GetCreateProject(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	c := createProjectDialog()
	return c.Render(r.Context(), w)
}

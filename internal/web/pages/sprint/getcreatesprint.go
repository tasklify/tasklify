package sprint

import (
	"net/http"
	"tasklify/internal/handlers"
)

func GetCreateSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	c := createSprintDialog()
	return c.Render(r.Context(), w)
}

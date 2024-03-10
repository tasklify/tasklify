package sprint

import (
	"net/http"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"
)

func GetSprint(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	c := createSprintDialog()
	return pages.Layout(c, "Tasklify").Render(r.Context(), w)
}

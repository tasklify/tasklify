package userstory

import (
	"net/http"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"
)

func GetUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	c := CreateUserStoryDialog()
	return pages.Layout(c, "Tasklify").Render(r.Context(), w)
}

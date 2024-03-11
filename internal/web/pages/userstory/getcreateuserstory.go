package userstory

import (
	"net/http"
	"tasklify/internal/handlers"
)

func GetUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	c := CreateUserStoryDialog()
	return c.Render(r.Context(), w)
}

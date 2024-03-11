package login

import (
	"net/http"
	"tasklify/internal/web/pages"
)

func GetLogin(w http.ResponseWriter, r *http.Request) error {
	c := login("Login")
	return pages.Layout(c, "Tasklify").Render(r.Context(), w)
}

package login

import (
	"net/http"
)

func GetLogin(w http.ResponseWriter, r *http.Request) error {
	c := login("Login")
	return c.Render(r.Context(), w)
}

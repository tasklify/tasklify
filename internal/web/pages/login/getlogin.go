package login

import (
	"net/http"
	"tasklify/internal/auth"
	"tasklify/internal/web/pages"
)

func GetLogin(w http.ResponseWriter, r *http.Request) error {
	// User already logged in
	_, err := auth.GetSession().GetUserID(r)
	if err == nil {
		w.Header().Set("HX-Redirect", "/dashboard")
		w.WriteHeader(http.StatusOK)
	}

	// User not logged in
	c := login()
	return pages.Layout(c, "Login", r).Render(r.Context(), w)
}

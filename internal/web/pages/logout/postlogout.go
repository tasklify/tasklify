package logout

import (
	"net/http"
	"tasklify/internal/auth"
)

func PostLogout(w http.ResponseWriter, r *http.Request) error {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return nil
	}

	err := auth.GetSession().Destroy(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
	return nil
}

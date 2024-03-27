package userSlug

import (
	"net/http"
	"strconv"
	"tasklify/internal/auth"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"

	"github.com/go-chi/chi/v5"
)

func DeleteUser(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	temp, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(err.Error())
		return c.Render(r.Context(), w)
	}
	userIDToDelete := uint(temp)

	err = auth.DeleteUser(params.UserID, userIDToDelete)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(err.Error())
		return c.Render(r.Context(), w)
	}

	w.Header().Set("HX-Redirect", "/users")
	w.WriteHeader(http.StatusOK)
	return nil
}

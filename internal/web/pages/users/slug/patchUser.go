package userSlug

import (
	"net/http"
	"strconv"
	"tasklify/internal/auth"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

type patchUserFormData struct {
	Username    string `schema:"username,required"`
	NewPassword string `schema:"new_password"`
	FirstName   string `schema:"first_name,required"`
	LastName    string `schema:"last_name,required"`
	Email       string `schema:"email,required"`
	SystemRole  string `schema:"system_role,required"`
	Password    string `schema:"password,required"`
}

var decoder = schema.NewDecoder()

func PatchUser(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	var patchUserFormData patchUserFormData
	err := decoder.Decode(&patchUserFormData, r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(err.Error())
		return c.Render(r.Context(), w)
	}

	temp, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(err.Error())
		return c.Render(r.Context(), w)
	}
	userIDToUpdate := uint(temp)

	var newPassword *string
	if len(patchUserFormData.NewPassword) != 0 {
		newPassword = &patchUserFormData.NewPassword
	}

	err = auth.UpdateUser(params.UserID,
		patchUserFormData.Password,
		userIDToUpdate,
		ptr.String(patchUserFormData.Username),
		newPassword,
		ptr.String(patchUserFormData.FirstName),
		ptr.String(patchUserFormData.LastName),
		ptr.String(patchUserFormData.Email),
		ptr.String(patchUserFormData.SystemRole),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(err.Error())
		return c.Render(r.Context(), w)
	}

	// w.Header().Set("HX-Redirect", "/users")
	w.WriteHeader(http.StatusOK)
	return nil
}

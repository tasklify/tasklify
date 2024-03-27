package users

import (
	"net/http"
	"tasklify/internal/auth"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"

	"github.com/aws/smithy-go/ptr"
	"github.com/gorilla/schema"
)

type postUsersFormData struct {
	Username   string `schema:"username,required"`
	Password   string `schema:"password,required"`
	PasswordRetype   string `schema:"password_retype,required"`
	FirstName  string `schema:"first_name,required"`
	LastName   string `schema:"last_name,required"`
	Email      string `schema:"email,required"`
	SystemRole string `schema:"system_role,required"`
}

var decoder = schema.NewDecoder()

func PostUsers(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	var postUsersFormData postUsersFormData
	err := decoder.Decode(&postUsersFormData, r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(err.Error())
		return c.Render(r.Context(), w)
	}

	err = auth.CreateUser(ptr.Uint(params.UserID), nil,
		postUsersFormData.Username,
		postUsersFormData.Password,
		postUsersFormData.PasswordRetype,
		postUsersFormData.FirstName,
		postUsersFormData.LastName,
		postUsersFormData.Email,
		postUsersFormData.SystemRole,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(err.Error())
		return c.Render(r.Context(), w)
	}

	w.Header().Set("HX-Redirect", "/users")
	w.WriteHeader(http.StatusOK)
	return nil
}

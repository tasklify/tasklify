package login

import (
	"net/http"
	"tasklify/internal/auth"
	"tasklify/internal/web/components/common"

	"github.com/gorilla/schema"
)

type loginFormData struct {
	Username string `schema:"username,required"`
	Password string `schema:"password,required"`
}

var decoder = schema.NewDecoder()

func PostLogin(w http.ResponseWriter, r *http.Request) error {
	var loginFormData loginFormData
	err := decoder.Decode(&loginFormData, r.PostForm)
	if err != nil {
		return err
	}

	userID, err := auth.AuthenticateUser(loginFormData.Username, loginFormData.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(err.Error())
		return c.Render(r.Context(), w)
	}

	err = auth.GetSession().Create(userID, w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError(err.Error())
		return c.Render(r.Context(), w)
	}

	w.Header().Set("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusOK)
	return nil
}

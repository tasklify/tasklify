package login

import (
	"log"
	"net/http"
	"tasklify/internal/auth"

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
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		c := loginError()
		return c.Render(r.Context(), w)
	}

	log.Printf("PostLogin: userID=%v", userID)

	if userID != 0 {
		err = auth.GetSession().Create(userID, w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			c := loginError()
			return c.Render(r.Context(), w)
		}

		w.Header().Set("HX-Redirect", "/dashboard")
		w.WriteHeader(http.StatusOK)
		return nil
	}

	w.WriteHeader(http.StatusUnauthorized)
	c := loginError()
	return c.Render(r.Context(), w)
}

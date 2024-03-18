package dashboard

import (
	"fmt"
	"net/http"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"
)

func Dashboard(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	myProjects, err := database.GetDatabase().GetUserProjects(params.UserID)
	if err != nil {
		return err
	}

	user, err := database.GetDatabase().GetUserByID(params.UserID)
	if err != nil {
		return err
	}

	c := pages.Index(fmt.Sprint(params.UserID), myProjects, user.SystemRole)
	return pages.Layout(c, "Tasklify").Render(r.Context(), w)
}

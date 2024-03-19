package userstory

import (
	"net/http"
	"strconv"
	"tasklify/internal/database"

	"tasklify/internal/handlers"
)

func GetUserStoryDetails(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	if err := r.ParseForm(); err != nil {
		return err
	}

	userStoryID, err := strconv.Atoi(r.FormValue("userStoryID"))
	if err != nil {
		return err
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(uint(userStoryID))
	if err != nil {
		return err
	}

	c := UserStoryDetailsDialog(*userStory)
	return c.Render(r.Context(), w)
}
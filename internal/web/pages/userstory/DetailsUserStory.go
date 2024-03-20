package userstory

import (
	"net/http"
	"tasklify/internal/database"

	"tasklify/internal/handlers"
)

func GetUserStoryDetails(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	type RequestData struct {
		UserStoryID uint `schema:"userStoryID,required"`
	}
	var requestData RequestData
	err := decoder.Decode(&requestData, r.URL.Query())

	if err != nil {
		return err
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(requestData.UserStoryID)
	if err != nil {
		return err
	}

	c := UserStoryDetailsDialog(*userStory)
	return c.Render(r.Context(), w)
}

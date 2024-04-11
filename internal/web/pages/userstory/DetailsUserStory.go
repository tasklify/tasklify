package userstory

import (
	"net/http"
	"strconv"
	"tasklify/internal/database"

	"github.com/go-chi/chi/v5"

	"tasklify/internal/handlers"
)

func GetUserStoryDetails(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		TabID *string `schema:"tabID"`
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	userStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))

	if err != nil {
		return err
	}

	var requestData RequestData
	err = decoder.Decode(&requestData, r.URL.Query())
	if err != nil {
		return err
	}

	var activeTab string
	if requestData.TabID == nil {
		activeTab = "details"
	} else {
		activeTab = *requestData.TabID
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(uint(userStoryID))
	if err != nil {
		return err
	}

	userStoryComments, err := database.GetDatabase().GetUserStoryComments(uint(userStoryID))
	if err != nil {
		return err
	}
	userStory.UserStoryComments = append(userStory.UserStoryComments, userStoryComments...)

	currentUser, err := database.GetDatabase().GetUserByID(params.UserID)
	if err != nil {
		return err
	}

	c := UserStoryDetailsDialog(*userStory, activeTab, *currentUser)
	return c.Render(r.Context(), w)

}

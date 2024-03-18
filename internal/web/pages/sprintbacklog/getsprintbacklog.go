package sprintbacklog

import (
	"net/http"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// GetSprintBacklog handles the request for fetching and displaying the sprint backlog.
func GetSprintBacklog(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		SprintID uint `schema:"sprintID,required"`
	}
	var requestData RequestData
	err := decoder.Decode(&requestData, r.URL.Query())
	if err != nil {
		return err
	}

	sprintID := requestData.SprintID

	//fetch sprint
	sprint := database.GetDatabase().GetSprintByID(sprintID)

	//get all users
	users, err := database.GetDatabase().GetUsers()
	if err != nil {
		return err
	}

	//map users to user ids
	userMap := mapUserstoUserIDs(users)

	c := sprintBacklog(sprintID, sprint, userMap)

	return pages.Layout(c, "Sprint Backlog").Render(r.Context(), w)
}

func mapUserstoUserIDs(users []database.User) (userMap map[uint]database.User) {
	userMap = make(map[uint]database.User)
	for _, user := range users {
		userMap[user.ID] = user
	}
	return
}
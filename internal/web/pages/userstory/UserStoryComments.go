package userstory

import (
	"net/http"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func PostComment(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		Body string `schema:"comment,required"`
	}

	userStoryIDInt, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}
	userStoryID := uint(userStoryIDInt)

	var req RequestData
	err = decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	newComment := database.UserStoryComment{
		UserStoryID: userStoryID,
		AuthorID:    params.UserID,
		Body:        req.Body,
	}

	err = database.GetDatabase().AddUserStoryComment(newComment)
	if err != nil {
		return err
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(userStoryID)
	if err != nil {
		return err
	}

	userStoryComments, err := database.GetDatabase().GetUserStoryComments(userStoryID)
	if err != nil {
		return err
	}
	userStory.UserStoryComments = append(userStory.UserStoryComments, userStoryComments...)

	currentUser, err := database.GetDatabase().GetUserByID(params.UserID)
	if err != nil {
		return err
	}

	c := UserStoryComments(*userStory, *currentUser)
	return c.Render(r.Context(), w)
}

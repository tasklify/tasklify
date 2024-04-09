package userstory

import (
	"net/http"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"

	"github.com/a-h/templ"
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

	c := UserStoryComments(*userStory, nil, *currentUser)
	return c.Render(r.Context(), w)
}

func GetEditComment(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	userStoryIDInt, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}
	userStoryID := uint(userStoryIDInt)

	commentIDInt, err := strconv.Atoi(chi.URLParam(r, "commentID"))
	if err != nil {
		return err
	}
	commentID := uint(commentIDInt)

	c, err := GetCommentList(userStoryID, &commentID, params.UserID)
	if err != nil {
		return err
	}

	return c.Render(r.Context(), w)
}

func GetCancelEditComment(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	userStoryIDInt, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}
	userStoryID := uint(userStoryIDInt)

	c, err := GetCommentList(userStoryID, nil, params.UserID)
	if err != nil {
		return err
	}

	return c.Render(r.Context(), w)
}

func PutComment(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		Body string `schema:"comment,required"`
	}

	userStoryIDInt, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}
	userStoryID := uint(userStoryIDInt)

	commentIDInt, err := strconv.Atoi(chi.URLParam(r, "commentID"))
	if err != nil {
		return err
	}
	commentID := uint(commentIDInt)

	var req RequestData
	err = decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	err = database.GetDatabase().EditUserStoryComment(userStoryID, commentID, req.Body)
	if err != nil {
		return err
	}

	c, err := GetCommentList(userStoryID, nil, params.UserID)
	if err != nil {
		return err
	}

	return c.Render(r.Context(), w)
}

func DeleteComment(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	userStoryIDInt, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}
	userStoryID := uint(userStoryIDInt)

	commentIDInt, err := strconv.Atoi(chi.URLParam(r, "commentID"))
	if err != nil {
		return err
	}
	commentID := uint(commentIDInt)

	err = database.GetDatabase().DeleteUserStoryComment(userStoryID, commentID)
	if err != nil {
		return err
	}

	c, err := GetCommentList(userStoryID, nil, params.UserID)
	if err != nil {
		return err
	}

	return c.Render(r.Context(), w)
}

func GetCommentList(userStoryID uint, inEditCommentID *uint, userID uint) (templ.Component, error) {
	userStory, err := database.GetDatabase().GetUserStoryByID(userStoryID)
	if err != nil {
		return nil, err
	}

	userStoryComments, err := database.GetDatabase().GetUserStoryComments(userStoryID)
	if err != nil {
		return nil, err
	}
	userStory.UserStoryComments = append(userStory.UserStoryComments, userStoryComments...)

	currentUser, err := database.GetDatabase().GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	var commentData *database.UserStoryComment
	if inEditCommentID != nil {
		commentData, err = database.GetDatabase().GetUserStoryCommentByID(*inEditCommentID)
		if err != nil {
			return nil, err
		}
	}

	return CommentList(*userStory, commentData, *currentUser), nil
}

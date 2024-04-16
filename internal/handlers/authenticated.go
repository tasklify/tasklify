package handlers

import (
	"context"
	"net/http"
	"tasklify/internal/auth"
)

type contextKeyUserID string

const (
	ContextKeyUserID contextKeyUserID = "user_id"
)

type AuthenticatedHandlerFunc func(w http.ResponseWriter, r *http.Request, params RequestParams) error

func (ahf AuthenticatedHandlerFunc) Serve(w http.ResponseWriter, r *http.Request) error {
	userID, err := auth.GetSession().GetUserID(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return nil
	}

	ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)

	return ahf(w, r.WithContext(ctx), RequestParams{UserID: userID})
}

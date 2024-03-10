package handlers

import (
	"errors"
	"net/http"
	"tasklify/internal/middlewares"
)

type AuthenticatedHandlerFunc func(w http.ResponseWriter, r *http.Request, params RequestParams) error

func (ahf AuthenticatedHandlerFunc) Serve(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value(middlewares.ContextKeyUserID).(uint)
	if !ok {
		http.Error(w, "Error retriving userID from context", http.StatusInternalServerError)
		return errors.New("error retriving userID from context")
	}

	return ahf(w, r, RequestParams{UserID: userID})
}

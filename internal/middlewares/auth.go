package middlewares

import (
	"context"
	"net/http"
	"tasklify/internal/auth"
)

type contextKeyUserID string

const (
	ContextKeyUserID contextKeyUserID = "user_id"
)

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetSession().GetUserID(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

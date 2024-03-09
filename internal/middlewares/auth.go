package middlewares

import (
	"context"
	"log"
	"net/http"
	"tasklify/internal/auth"
)

type contextKeyUserId string

const (
	ContextKeyUserId contextKeyUserId = "user_id"
)

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := auth.GetSession().GetUserId(r)
		if err != nil {
			log.Printf("Middleware: AuthUser: %v\n", err)

			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserId, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

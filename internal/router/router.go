package router

import (
	"net/http"
	"tasklify/internal/handlers"
	"tasklify/internal/middlewares"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			middlewares.TextHTMLMiddleware,
			middlewares.CSPMiddleware,
			// TODO: https://github.com/gorilla/csrf
			middleware.Compress(5),
		)

		// Public
		r.NotFound(handlers.NewNotFoundHandler().ServeHTTP)
		r.Get("/", handlers.NewHomeHandler().ServeHTTP)
		r.Get("/about", handlers.NewAboutHandler().ServeHTTP)
		r.Get("/register", handlers.NewGetRegisterHandler().ServeHTTP)
		r.Post("/register", handlers.NewPostRegisterHandler(handlers.PostRegisterHandlerParams{
			// UserStore: userStore,
		}).ServeHTTP)
		r.Get("/login", handlers.NewGetLoginHandler().ServeHTTP)
		r.Post("/login", handlers.NewPostLoginHandler(handlers.PostLoginHandlerParams{
			// UserStore: userStore,
			// TokenAuth: tokenAuth,
		}).ServeHTTP)

		// Secure
		r.Group(func(r chi.Router) {
			r.Use(
				middlewares.AuthUser,
			)

			r.Get("/dashboard", handlers.NewHomeHandler().ServeHTTP)
		})
	})

	return r
}

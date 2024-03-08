package router

import (
	"net/http"
	"tasklify/internal/middlewares"
	"tasklify/internal/pages"
	"tasklify/internal/pages/about"
	"tasklify/internal/pages/login"

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
		r.NotFound(pages.NewNotFoundHandler().ServeHTTP)
		r.Get("/", pages.NewHomeHandler().ServeHTTP)
		r.Get("/about", about.NewAboutHandler().ServeHTTP)
		r.Get("/login", login.NewGetLoginHandler().ServeHTTP)
		r.Post("/login", login.NewPostLoginHandler(login.PostLoginHandlerParams{
			// UserStore: userStore,
			// TokenAuth: tokenAuth,
		}).ServeHTTP)

		// Secure
		r.Group(func(r chi.Router) {
			r.Use(
				middlewares.AuthUser,
			)

			r.Get("/dashboard", pages.NewHomeHandler().ServeHTTP)
		})
	})

	return r
}

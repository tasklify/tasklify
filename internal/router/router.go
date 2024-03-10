package router

import (
	"net/http"
	"tasklify/internal/handlers"
	"tasklify/internal/middlewares"
	"tasklify/internal/web/pages"
	"tasklify/internal/web/pages/about"
	"tasklify/internal/web/pages/dashboard"
	"tasklify/internal/web/pages/login"
	"tasklify/internal/web/pages/project"
	"tasklify/internal/web/pages/sprint"
	"tasklify/internal/web/pages/userstory"

	ghandlers "github.com/gorilla/handlers"

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
			// TODO: CORS
			middleware.Compress(5),
		)

		// Public

		r.Handle("/", ghandlers.MethodHandler{
			"GET": handlers.UnifiedHandler(handlers.PlainHandlerFunc(pages.Home)),
		})
		r.Handle("/about", ghandlers.MethodHandler{
			"GET": handlers.UnifiedHandler(handlers.PlainHandlerFunc(about.About)),
		})
		r.Handle("/login", ghandlers.MethodHandler{
			"GET":  handlers.UnifiedHandler(handlers.PlainHandlerFunc(login.GetLogin)),
			"POST": handlers.UnifiedHandler(handlers.PlainHandlerFunc(login.PostLogin)),
		})

		r.NotFound(handlers.UnifiedHandler(handlers.PlainHandlerFunc(pages.NotFound)))

		// Secure
		r.Group(func(r chi.Router) {
			r.Use(
				middlewares.AuthUser,
			)

			r.Handle("/dashboard", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(dashboard.Dashboard)),
			})
			r.Handle("/project", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.GetCreateProject)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.PostProject)),
			})
			r.Handle("/sprint", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprint.GetSprint)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprint.PostSprint)),
			})
			r.Handle("/createuserstory", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.GetUserStory)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.PostUserStory)),
			})
		})
	})

	return r
}

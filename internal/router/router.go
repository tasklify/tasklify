package router

import (
	"net/http"
	"tasklify/internal/handlers"
	"tasklify/internal/middlewares"
	"tasklify/internal/web/pages"
	"tasklify/internal/web/pages/about"
	"tasklify/internal/web/pages/login"
	"tasklify/internal/web/pages/logout"
	"tasklify/internal/web/pages/productbacklog"
	"tasklify/internal/web/pages/project"
	"tasklify/internal/web/pages/sprint"
	"tasklify/internal/web/pages/sprintbacklog"
	"tasklify/internal/web/pages/task"
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
			r.Handle("/logout", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.PlainHandlerFunc(logout.PostLogout)),
			})
			r.Handle("/create-project", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.GetCreateProject)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.PostCreateProject)),
			})
			r.Handle("/project-developer", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.PostAddProjectDeveloper)),
			})
			r.Handle("/remove-project-developer", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.RemoveProjectDeveloper)),
			})
			r.Handle("/createsprint", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprint.GetCreateSprint)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprint.PostSprint)),
			})
			r.Handle("/createuserstory", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.GetUserStory)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.PostUserStory)),
			})
			r.Handle("/create-task", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.GetCreateTask)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.PostTask)),
			})
			r.Handle("/productbacklog", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.GetProductBacklog)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.PostAddUserStoryToSprint)),
			})
			r.Handle("/userstory/remove-from-sprint", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.RemoveUserStoryFromSprint)),
			})
			r.Handle("/userstory/details", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.GetUserStoryDetails)),
			})
			r.Handle("/sprintbacklog", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprintbacklog.GetSprintBacklog)),
			})
			r.Handle("/userstory/accept", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.PostUserStoryAccepted)),
			})
			r.Handle("/userstory/reject", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.PostUserStoryRejected)),
			})
			r.Handle("/userstory/rejectioncomment", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.PostRejectionComment)),
			})
		})
	})

	return r
}

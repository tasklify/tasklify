package router

import (
	"net/http"
	"tasklify/internal/handlers"
	"tasklify/internal/middlewares"
	"tasklify/internal/web/pages"
	"tasklify/internal/web/pages/about"
	"tasklify/internal/web/pages/dashboard"
	"tasklify/internal/web/pages/login"
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

			r.Handle("/dashboard", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(dashboard.Dashboard)),
			})

			// ===== Create Project endpoints =====
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

			// ===== Edit Project endpoints =====
			r.Handle("/edit-project-info/{projectID}", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.GetEditProjectInfo)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.UpdateProjectInfo)),
			})
			r.Handle("/edit-project-members/{projectID}", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.GetEditProjectMembers)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.UpdateProjectMembers)),
			})
			r.Handle("/edit-project-members/{projectID}/remove-project-developer", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.EditProjectRemoveDeveloper)),
			})
			r.Handle("/edit-project-members/{projectID}/add-project-developer", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.EditProjectAddDeveloper)),
			})
			r.Handle("/edit-project-members/{projectID}/change-product-owner", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.EditProjectChangeOwner)),
			})
			r.Handle("/edit-project-members/{projectID}/change-scrum-master", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(project.EditProjectChangeMaster)),
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

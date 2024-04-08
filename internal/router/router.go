package router

import (
	"net/http"
	"tasklify/internal/docs"
	"tasklify/internal/handlers"
	"tasklify/internal/middlewares"
	"tasklify/internal/web/pages"
	"tasklify/internal/web/pages/dashboard"
	"tasklify/internal/web/pages/login"
	"tasklify/internal/web/pages/logout"
	"tasklify/internal/web/pages/productbacklog"
	"tasklify/internal/web/pages/project"
	"tasklify/internal/web/pages/projectinfo"
	"tasklify/internal/web/pages/projectwall"
	"tasklify/internal/web/pages/sprint"
	"tasklify/internal/web/pages/sprintbacklog"
	"tasklify/internal/web/pages/task"
	"tasklify/internal/web/pages/users"
	userSlug "tasklify/internal/web/pages/users/slug"
	"tasklify/internal/web/pages/userstory"

	ghandlers "github.com/gorilla/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.Get("/docs", handlers.DocsHandler(docs.RenderDoc))
	r.Handle("/docs/*", http.StripPrefix("/docs/", handlers.DocsHandler(docs.RenderDoc)))

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
			"GET": handlers.UnifiedHandler(handlers.PlainHandlerFunc(pages.Index)),
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

			// ===== Users endpoints =====
			r.Handle("/users", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(users.Users)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(users.PostUsers)),
			})
			r.Handle("/users/new", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(users.GetNewUser)),
			})
			r.Handle("/users/{userID}", ghandlers.MethodHandler{
				"GET":    handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userSlug.GetUser)),
				"PATCH":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userSlug.PatchUser)),
				"DELETE": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userSlug.DeleteUser)),
			})
			r.Handle("/users/{userID}/delete", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userSlug.GetDeleteUser)),
			})

			// ===== Create Project endpoints =====
			r.Handle("/dashboard", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(dashboard.Dashboard)),
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

			r.Handle("/project-info/{projectID}", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(projectinfo.GetProjectInfo)),
			})

			r.Handle("/project-wall/{projectID}", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(projectwall.GetProjectWall)),
			})
			r.Handle("/project-wall/{projectID}/post", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(projectwall.GetNewPostDialog)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(projectwall.AddNewPost)),
			})
			r.Handle("/project-wall/{projectID}/post/{postID}", ghandlers.MethodHandler{
				"GET":    handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(projectwall.GetEditPost)),
				"PUT":    handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(projectwall.PutPost)),
				"DELETE": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(projectwall.DeletePost)),
			})

			r.Handle("/{projectID}/createsprint", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprint.GetCreateSprint)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprint.PostSprint)),
			})
			r.Handle("/{projectID}/createuserstory", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.GetUserStory)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.PostUserStory)),
			})
			r.Handle("/userstory/{userStoryID}/create-task", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.GetCreateTask)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.PostTask)),
			})
			r.Handle("/productbacklog", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.GetProductBacklog)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.PostAddUserStoryToSprint)),
			})
			r.Handle("/userstory/{userStoryID}/remove-from-sprint", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.RemoveUserStoryFromSprint)),
			})
			r.Handle("/userstory/{userStoryID}/details", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.GetUserStoryDetails)),
			})
			r.Handle("/userstory/{userStoryID}/accept", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.PostUserStoryAccepted)),
			})
			r.Handle("/userstory/{userStoryID}/reject", ghandlers.MethodHandler{
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.GetUserStoryRejected)),
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(productbacklog.PostUserStoryRejected)),
			})
			r.Handle("/userstory/{userStoryID}/comment", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.PostComment)),
			})
			r.Handle("/userstory/{userStoryID}/comment/{commentID}", ghandlers.MethodHandler{
				"GET":    handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.GetEditComment)),
				"PUT":    handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.PutComment)),
				"DELETE": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.DeleteComment)),
			})
			r.Handle("/userstory/{userStoryID}/comment/{commentID}/cancel-edit", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(userstory.GetCancelEditComment)),
			})
			r.Handle("/sprintbacklog/{sprintID}", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprintbacklog.GetSprintBacklog)),
			})
			r.Handle("/sprintbacklog/{sprintID}/task/{taskID}/details", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.GetTaskDetails)),
			})
			r.Handle("/sprintbacklog/{sprintID}/task/{taskID}/unassign", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprintbacklog.UnassignTask)),
			})
			r.Handle("/sprintbacklog/{sprintID}/task/{taskID}/assign", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(sprintbacklog.AssignTask)),
			})
			r.Handle("/sprintbacklog/{sprintID}/task/{taskID}/start", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.StartWorkSession)),
			})
			r.Handle("/sprintbacklog/{sprintID}/task/{taskID}/logged", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.GetLoggedTime)),
			})
			r.Handle("/sprintbacklog/{sprintID}/task/{taskID}/logged/session/{workSessionID}/stop", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.StopWorkSession)),
			})
			r.Handle("/sprintbacklog/{sprintID}/task/{taskID}/logged/session/{workSessionID}/change", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.PostChangeDuration)),
				"GET":  handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.GetChangeDuration)),
			})
			r.Handle("/sprintbacklog/{sprintID}/task/{taskID}/logged/session/{workSessionID}/resume", ghandlers.MethodHandler{
				"POST": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.ResumeWorkSession)),
			})
			r.Handle("/sprintbacklog/{sprintID}/task/{taskID}/logged/session/{workSessionID}/unfinished", ghandlers.MethodHandler{
				"GET": handlers.UnifiedHandler(handlers.AuthenticatedHandlerFunc(task.GetUnfinishedSessionDialog)),
			})

		})
	})

	return r
}

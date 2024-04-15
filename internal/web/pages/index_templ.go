// Code generated by templ - DO NOT EDIT.

package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"net/http"
	"tasklify/internal/auth"
	"tasklify/internal/database"
)

func Index(w http.ResponseWriter, r *http.Request) error {
	sessionManager := auth.GetSession()
	userID, err := sessionManager.GetUserID(r)
	if err != nil {
		// User is not logged in; show the guest index
		c := guestIndex()
		return Layout(c, "Tasklify", r).Render(r.Context(), w)
	}

	myProjects, err := database.GetDatabase().GetUserProjects(userID)
	if err != nil {
		// Handle error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	user, err := database.GetDatabase().GetUserByID(userID)
	if err != nil {
		// Handle error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	c := index(fmt.Sprint(userID), myProjects, user.SystemRole)
	return Layout(c, "Tasklify", r).Render(r.Context(), w)
}

func index(userID string, myProjects []database.Project, user_SystemRole database.SystemRole) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-base-200 min-h-screen\"><div class=\"p-10\"><div class=\"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6\"><!-- Add New Project -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if auth.GetAuthorization().HasSystemPermission(user_SystemRole, "/project", auth.ActionCreate) == nil {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card card-compact bg-base-100 shadow-xl hover:shadow-2xl transition-shadow flex justify-center items-center\"><div class=\"card-body\"><div><h2 class=\"card-title\">Add new project</h2><a href=\"/docs/projects#creating-projects\" target=\"_blank\" class=\"help-button\">?</a></div><div class=\"flex justify-center\"><button class=\"btn btn-primary btn-circle btn-lg\" hx-get=\"/create-project\" hx-target=\"#dialog\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-6 w-6\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\" stroke-width=\"2\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M12 4v16m8-8H4\"></path></svg></button></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		for _, project := range myProjects {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card card-compact bg-base-100 shadow-xl hover:shadow-2xl transition-shadow\"><div class=\"card-body\"><h2 class=\"card-title\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(project.Title)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/index.templ`, Line: 61, Col: 45}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2><p class=\"whitespace-break-spaces\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(project.Description)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/index.templ`, Line: 62, Col: 63}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><div class=\"card-actions justify-end\"><form hx-get=\"/productbacklog\" hx-target=\"#whole\" hx-swap=\"innerHTML\" hx-push-url=\"true\"><input type=\"hidden\" id=\"projectID\" name=\"projectID\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprint(project.ID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <button class=\"btn btn-primary\">View project</button></form></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div id=\"project-dialog\"></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func guestIndex() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"hero min-h-screen bg-base-200\"><div class=\"hero-content flex-col lg:flex-row-reverse\"><div class=\"text-center lg:text-left\"><h1 class=\"text-5xl font-bold\">Login now!</h1></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

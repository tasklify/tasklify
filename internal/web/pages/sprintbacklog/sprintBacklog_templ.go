// Code generated by templ - DO NOT EDIT.

package sprintbacklog

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "tasklify/internal/database"
import "fmt"
import "strconv"
import "slices"
import "time"
import "tasklify/internal/web/components/common"
import "sort"

// sprintBacklog.templ:
func sprintBacklog(sprint *database.Sprint, projectRoles []database.ProjectRole, userId uint) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-base-100 min-h-screen\"><div class=\"container mx-auto p-2\"><!-- Back to Product Backlog Button --><div class=\"mb-2\"><div class=\"flex justify-between items-center\"><form hx-get=\"/productbacklog\" hx-target=\"#whole\" hx-push-url=\"true\"><input type=\"hidden\" id=\"projectID\" name=\"projectID\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprint(sprint.ProjectID)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <button class=\"btn btn-primary btn-xs\">Back To Product Backlog</button></form><a href=\"https://github.com/tasklify/tasklify/tree/main?tab=readme-ov-file#sprint-backlog\" target=\"_blank\" class=\"help-button ml-4\" style=\"position: relative;top:0px;right:0px\">?</a></div><!-- Sprint Details --><div class=\"sprint-details flex justify-between items-center mb-4 p-2 shadow-lg bg-white rounded-lg\"><h2 class=\"font-bold text-lg\" style=\"padding-left: 15px;\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(sprint.Title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 25, Col: 77}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2><p class=\"text-sm\" style=\"padding-right: 10px;\">(")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(sprint.StartDate.Format("Mon Jan _2 2006"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 26, Col: 98}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" - ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(sprint.EndDate.Format("Mon Jan _2 2006"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 26, Col: 145}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(")</p></div></div><!-- User Stories -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, story := range sprint.UserStories {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"border border-base-300 mb-2 rounded-lg\"><div class=\"flex justify-between items-center bg-base-300 p-2\"><div><div class=\"font-medium text-lg\" style=\"padding-left: 15px;\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(story.Title)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 34, Col: 81}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div class=\"flex gap-2 mt-2\"><button title=\"Details\" hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/userstory/%v/details", story.ID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#dialog\" class=\"btn btn-xs\"><svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"w-6 h-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M6.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM12.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM18.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Z\"></path></svg></button><!-- Create task Button -->")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if *story.Realized {
				if slices.Contains(projectRoles, database.ProjectRoleMaster) || slices.Contains(projectRoles, database.ProjectRoleDeveloper) {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form hx-get=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/userstory/%v/create-task", story.ID)))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#dialog\" title=\"Can&#39;t add tasks to realized user story\"><input type=\"hidden\" id=\"sprintID\" name=\"sprintID\" value=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(strconv.Itoa(int(sprint.ID))))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <input type=\"hidden\" id=\"projectID\" name=\"projectID\" value=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(strconv.Itoa(int(story.ProjectID))))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <button class=\"btn btn-xs btn-primary\" style=\"margin-right: 10px;\" disabled><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"white\" class=\"w-5 h-5\"><path d=\"M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z\"></path></svg></button></form>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			} else {
				if slices.Contains(projectRoles, database.ProjectRoleMaster) || slices.Contains(projectRoles, database.ProjectRoleDeveloper) {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form hx-get=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/userstory/%v/create-task", story.ID)))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#dialog\"><input type=\"hidden\" id=\"sprintID\" name=\"sprintID\" value=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(strconv.Itoa(int(sprint.ID))))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <input type=\"hidden\" id=\"projectID\" name=\"projectID\" value=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(strconv.Itoa(int(story.ProjectID))))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <button title=\"Create task\" class=\"btn btn-xs btn-primary\" style=\"margin-right: 10px;\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"white\" class=\"w-5 h-5\"><path d=\"M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z\"></path></svg></button></form>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div class=\"collapse-body card-body bg-base-100 p-2\"><table class=\"w-full\"><thead><tr class=\"text-sm font-bold\"><th class=\"text-left p-2\" style=\"width: 35%;\"></th><th class=\"text-left p-2\" style=\"width: 20%;\">Assigned To</th><th class=\"text-right p-2\" style=\"width: 8%;\">Estimated</th><th class=\"text-right p-2\" style=\"width: 8%;\">Logged</th><th class=\"text-right p-2\" style=\"width: 8%;\">Remaining</th><th class=\"text-center p-2\" style=\"width: 15%;\">Status</th><th class=\"text-right p-2\" style=\"width: 6%;\"></th></tr></thead> <tbody>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, task := range story.Tasks {
				if !*task.UserAccepted && task.UserID != nil && *task.UserID == userId {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<tr class=\"bg-warning shadow-lg\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					templ_7745c5c3_Err = TaskRowSmall(task, userId, sprint.ID).Render(ctx, templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</tr>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<tr class=\"bg-base-300 shadow-lg\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					templ_7745c5c3_Err = TaskRowSmall(task, userId, sprint.ID).Render(ctx, templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</tr>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</tbody><tr class=\"text-left text-sm font-bold\"><td class=\"p-2\" colspan=\"2\">Total:</td><td class=\"p-2 text-right\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(common.FormatDuration(GetSumOfTimeEstimates(story.Tasks)))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 98, Col: 94}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</td><td class=\"p-2 text-right\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 string
			templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(common.FormatDuration(GetSumOfLogs(story.Tasks)))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 99, Col: 85}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</td><td class=\"p-2 text-right\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 string
			templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(common.FormatDuration(GetSumOfRemaining(story.Tasks)))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 100, Col: 90}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</td><td class=\"p-2\" colspan=\"2\"></td></tr></table></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func TaskRowSmall(task database.Task, userId uint, sprintID uint) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var9 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var9 == nil {
			templ_7745c5c3_Var9 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- Flex container for aligning assignee and details button --><td class=\"font-bold text-md pl-4\"><div class=\"flex items-center\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !*task.UserAccepted && task.UserID != nil && *task.UserID == userId {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"mr-2\" title=\"The task needs to be accepted or rejected.\"><svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"w-6 h-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z\"></path></svg></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(*task.Title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 122, Col: 16}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></td><td class=\"text-left p-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if task.UserID == nil {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex items-center\"><span class=\"badge badge-ghost mr-2\">Unassigned</span> <span><button hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v/assign", sprintID, task.ID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"btn btn-xs\" title=\"Assign\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"currentColor\" class=\"w-5 h-5\"><path d=\"M10 5a3 3 0 1 1-6 0 3 3 0 0 1 6 0ZM1.615 16.428a1.224 1.224 0 0 1-.569-1.175 6.002 6.002 0 0 1 11.908 0c.058.467-.172.92-.57 1.174A9.953 9.953 0 0 1 7 18a9.953 9.953 0 0 1-5.385-1.572ZM16.25 5.75a.75.75 0 0 0-1.5 0v2h-2a.75.75 0 0 0 0 1.5h2v2a.75.75 0 0 0 1.5 0v-2h2a.75.75 0 0 0 0-1.5h-2v-2Z\"></path></svg></button></span></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if !*task.UserAccepted {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!--accept in decline se vidi samo če smo current user--> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if *task.UserID == userId {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex items-center\"><span class=\"badge bg-primary mr-2\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var11 string
				templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(GetUserFirstAndLastNameFromID(*task.UserID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 141, Col: 86}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> <button hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v/assign", sprintID, task.ID)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"btn btn-xs mr-1\" title=\"Accept task\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"currentColor\" class=\"w-5 h-5\"><path fill-rule=\"evenodd\" d=\"M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5Z\" clip-rule=\"evenodd\"></path></svg></button> <button hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v/unassign", sprintID, task.ID)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"btn btn-xs\" title=\"Decline task\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"currentColor\" class=\"w-5 h-5\"><path fill-rule=\"evenodd\" d=\"M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16ZM8.28 7.22a.75.75 0 0 0-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 1 0 1.06 1.06L10 11.06l1.72 1.72a.75.75 0 1 0 1.06-1.06L11.06 10l1.72-1.72a.75.75 0 0 0-1.06-1.06L10 8.94 8.28 7.22Z\" clip-rule=\"evenodd\"></path></svg></button></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"badge bg-primary\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var12 string
				templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(GetUserFirstAndLastNameFromID(*task.UserID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 154, Col: 80}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> <span class=\"text-xs italic\">(pending)</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex items-center\"><span class=\"badge bg-primary mr-2\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var13 string
			templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(GetUserFirstAndLastNameFromID(*task.UserID))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 159, Col: 85}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if *task.UserID == userId {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v/unassign", sprintID, task.ID)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"btn btn-xs\" title=\"Unassign\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" class=\"w-5 h-5\"><path d=\"M11 5a3 3 0 1 1-6 0 3 3 0 0 1 6 0ZM2.046 15.253c-.058.468.172.92.57 1.175A9.953 9.953 0 0 0 8 18c1.982 0 3.83-.578 5.384-1.573.398-.254.628-.707.57-1.175a6.001 6.001 0 0 0-11.908 0ZM12.75 7.75a.75.75 0 0 0 0 1.5h5.5a.75.75 0 0 0 0-1.5h-5.5Z\"></path></svg></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</td><td class=\"text-right p-2\"><span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var14 string
		templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(common.FormatDuration(task.TimeEstimate))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 171, Col: 50}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></td><td class=\"text-right p-2\"><span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var15 string
		templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(common.FormatDuration(GetTotalTimeSpentOnTask(task)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 174, Col: 62}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></td><td class=\"text-right p-2\"><span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var16 string
		templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(common.FormatDuration(GetRemainingTimeOnTask(task)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 177, Col: 61}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></td><td class=\"text-center p-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if GetTaskStatus(task) == "Unassigned" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"text-sm badge bg-unassigned-300\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var17 string
			templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs(GetTaskStatus(task))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 181, Col: 70}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var17))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if GetTaskStatus(task) == "Pending" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"text-sm badge bg-pending-300\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var18 string
			templ_7745c5c3_Var18, templ_7745c5c3_Err = templ.JoinStringErrs(GetTaskStatus(task))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 183, Col: 67}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var18))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if GetTaskStatus(task) == "Assigned" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"text-sm badge bg-assigned-300\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var19 string
			templ_7745c5c3_Var19, templ_7745c5c3_Err = templ.JoinStringErrs(GetTaskStatus(task))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 185, Col: 68}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var19))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if GetTaskStatus(task) == "Active" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"text-sm badge bg-active-300\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var20 string
			templ_7745c5c3_Var20, templ_7745c5c3_Err = templ.JoinStringErrs(GetTaskStatus(task))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 187, Col: 66}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var20))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if GetTaskStatus(task) == "Done" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"text-sm badge bg-done-300\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var21 string
			templ_7745c5c3_Var21, templ_7745c5c3_Err = templ.JoinStringErrs(GetTaskStatus(task))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/web/pages/sprintbacklog/sprintBacklog.templ`, Line: 189, Col: 64}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var21))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</td><td class=\"text-right pr-2 rounded-r-lg\" style=\"margin-top: 10px; margin-bottom: -10px; display: flex; justify-content: flex-end;\"><div class=\"start-button\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if (GetTaskStatus(task) == "Assigned") && (userId == *task.UserID) {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form hx-post=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v/start", sprintID, task.ID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><button type=\"submit\" class=\"btn btn-xs\">Start</button></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"logged-button\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if (GetTaskStatus(task) == "Active") && (userId == *task.UserID) {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v/sessions", sprintID, task.ID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#dialog\"><button class=\"btn btn-xs\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"currentColor\" class=\"w-5 h-5\"><path fill-rule=\"evenodd\" d=\"M4.5 2A1.5 1.5 0 0 0 3 3.5v13A1.5 1.5 0 0 0 4.5 18h11a1.5 1.5 0 0 0 1.5-1.5v-13A1.5 1.5 0 0 0 15.5 2h-11zm1 3a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 0 1h-8a.5.5 0 0 1-.5-.5zm0 3a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 0 1h-8a.5.5 0 0 1-.5-.5zm0 3a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 0 1h-8a.5.5 0 0 1-.5-.5zm0 3a.5.5 0 0 1 .5-.5h5a.5.5 0 0 1 0 1h-5a.5.5 0 0 1-.5-.5z\" clip-rule=\"evenodd\"></path></svg></button></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><!-- Details task Button --><span title=\"Details\"><div class=\"details-button ml-2\"><form hx-get=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v/details", sprintID, task.ID)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#dialog\"><button class=\"btn btn-xs\"><svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"w-6 h-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M6.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM12.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM18.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Z\"></path></svg></button></form></div></span><!-- Edit task Button --><span title=\"Edit task\"><div class=\"ml-2\"><button hx-get=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v", sprintID, task.ID)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#dialog\" class=\"btn btn-xs bg-primary\" style=\"position: relative; z-index: 1;\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"white\" class=\"w-5 h-5\"><path d=\"m5.433 13.917 1.262-3.155A4 4 0 0 1 7.58 9.42l6.92-6.918a2.121 2.121 0 0 1 3 3l-6.92 6.918c-.383.383-.84.685-1.343.886l-3.154 1.262a.5.5 0 0 1-.65-.65Z\"></path> <path d=\"M3.5 5.75c0-.69.56-1.25 1.25-1.25H10A.75.75 0 0 0 10 3H4.75A2.75 2.75 0 0 0 2 5.75v9.5A2.75 2.75 0 0 0 4.75 18h9.5A2.75 2.75 0 0 0 17 15.25V10a.75.75 0 0 0-1.5 0v5.25c0 .69-.56 1.25-1.25 1.25h-9.5c-.69 0-1.25-.56-1.25-1.25v-9.5Z\"></path></svg></button></div></span><!-- Delete task Button -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if GetTaskStatus(task) == "Unassigned" && GetTotalTimeSpentOnTask(task) == 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span title=\"Delete task\"><div class=\"ml-2\"><button hx-delete=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v", sprintID, task.ID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#dialog\" class=\"btn btn-xs bg-primary\" style=\"position: relative; z-index: 1;\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"white\" class=\"w-5 h-5\"><path fill-rule=\"evenodd\" d=\"M8.75 1A2.75 2.75 0 0 0 6 3.75v.443c-.795.077-1.584.176-2.365.298a.75.75 0 1 0 .23 1.482l.149-.022.841 10.518A2.75 2.75 0 0 0 7.596 19h4.807a2.75 2.75 0 0 0 2.742-2.53l.841-10.52.149.023a.75.75 0 0 0 .23-1.482A41.03 41.03 0 0 0 14 4.193V3.75A2.75 2.75 0 0 0 11.25 1h-2.5ZM10 4c.84 0 1.673.025 2.5.075V3.75c0-.69-.56-1.25-1.25-1.25h-2.5c-.69 0-1.25.56-1.25 1.25v.325C8.327 4.025 9.16 4 10 4ZM8.58 7.72a.75.75 0 0 0-1.5.06l.3 7.5a.75.75 0 1 0 1.5-.06l-.3-7.5Zm4.34.06a.75.75 0 1 0-1.5-.06l-.3 7.5a.75.75 0 1 0 1.5.06l.3-7.5Z\" clip-rule=\"evenodd\"></path></svg></button></div></span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span title=\"Assigned tasks and tasks with logged time can not be deleted.\"><div class=\"ml-2\"><button hx-delete=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/sprintbacklog/%v/task/%v", sprintID, task.ID)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#dialog\" class=\"btn btn-xs bg-primary\" style=\"position: relative; z-index: 1;\" disabled><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"white\" class=\"w-5 h-5\"><path fill-rule=\"evenodd\" d=\"M8.75 1A2.75 2.75 0 0 0 6 3.75v.443c-.795.077-1.584.176-2.365.298a.75.75 0 1 0 .23 1.482l.149-.022.841 10.518A2.75 2.75 0 0 0 7.596 19h4.807a2.75 2.75 0 0 0 2.742-2.53l.841-10.52.149.023a.75.75 0 0 0 .23-1.482A41.03 41.03 0 0 0 14 4.193V3.75A2.75 2.75 0 0 0 11.25 1h-2.5ZM10 4c.84 0 1.673.025 2.5.075V3.75c0-.69-.56-1.25-1.25-1.25h-2.5c-.69 0-1.25.56-1.25 1.25v.325C8.327 4.025 9.16 4 10 4ZM8.58 7.72a.75.75 0 0 0-1.5.06l.3 7.5a.75.75 0 1 0 1.5-.06l-.3-7.5Zm4.34.06a.75.75 0 1 0-1.5-.06l-.3 7.5a.75.75 0 1 0 1.5.06l.3-7.5Z\" clip-rule=\"evenodd\"></path></svg></button></div></span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</td><tr><td style=\"height: 10px;\"></td></tr>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func GetSumOfTimeEstimates(tasks []database.Task) (sum time.Duration) {
	sum = 0
	for _, task := range tasks {
		if task.TimeEstimate != 0 {
			sum += task.TimeEstimate
		}
	}
	return
}

func GetSumOfLogs(tasks []database.Task) (sum time.Duration) {
	sum = 0
	for _, task := range tasks {
		sum += GetTotalTimeSpentOnTask(task)
	}
	return sum
}

func GetSumOfRemaining(tasks []database.Task) (sum time.Duration) {
	sum = 0
	for _, task := range tasks {
		sum += GetRemainingTimeOnTask(task)
	}
	return sum
}

func GetTotalTimeSpentOnTask(task database.Task) time.Duration {
	var totalTimeSpent time.Duration = 0
	for _, session := range task.WorkSessions {
		totalTimeSpent += session.Duration
	}
	return totalTimeSpent
}

func GetRemainingTimeOnTask(task database.Task) time.Duration {
	sessions, err := database.GetDatabase().GetWorkSessionsForTask(task.ID)
	if err != nil {
		return task.TimeEstimate
	}

	sessions = sortWorkSessionsByDate(sessions)
	if len(sessions) == 0 {
		return task.TimeEstimate
	}

	lastSession := sessions[len(sessions)-1]
	return lastSession.Remaining
}

func sortWorkSessionsByDate(workSessions []database.WorkSession) []database.WorkSession {
	sort.Slice(workSessions, func(i, j int) bool {
		return workSessions[i].StartTime.Before(workSessions[j].StartTime)
	})

	return workSessions
}

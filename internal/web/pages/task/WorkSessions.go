package task

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"
	"time"

	"github.com/go-chi/chi/v5"
)

func sortWorkSessionsByDate(workSessions []database.WorkSession) []database.WorkSession {
	sort.Slice(workSessions, func(i, j int) bool {
		return workSessions[i].StartTime.Before(workSessions[j].StartTime)
	})

	return workSessions
}

func StartWorkSession(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}

	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	if *task.Status != database.StatusInProgress {
		*task.Status = database.StatusInProgress
		err = database.GetDatabase().UpdateTask(task)
		if err != nil {
			return err
		}
		err = database.GetDatabase().CreateWorkSession(params.UserID, uint(taskID))
		if err != nil {
			return err
		}
		w.Header().Set("HX-Redirect", fmt.Sprint("/sprintbacklog/", sprintID))
		w.WriteHeader(http.StatusSeeOther)

		return nil
	}

	err = database.GetDatabase().CreateWorkSession(params.UserID, uint(taskID))
	if err != nil {
		return err
	}

	todaysWS, otherWS, err := FetchSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task)
	return c.Render(r.Context(), w)
}

func ResumeWorkSession(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	workSessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	session, err := database.GetDatabase().GetWorkSessionByID(uint(workSessionID))
	if err != nil {
		return err
	}

	session.StartTime = time.Now()
	session.EndTime = nil
	err = database.GetDatabase().UpdateWorkSession(session)
	if err != nil {
		return err
	}

	todaysWS, otherWS, err := FetchSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task)
	return c.Render(r.Context(), w)
}

func StopWorkSession(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	workSessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}
	
	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	workSession, err := database.GetDatabase().GetWorkSessionByID(uint(workSessionID))
	if err != nil {
		return err
	}

	endTime := time.Now()
	workSession.Remaining = workSession.Remaining - endTime.Sub(workSession.StartTime)
	if workSession.Remaining < 0 {
		workSession.Remaining = 0
	}
	workSession.EndTime = &endTime
	workSession.Duration += endTime.Sub(workSession.StartTime)
	err = database.GetDatabase().UpdateWorkSession(workSession)
	if err != nil {
		return err
	}
	
	todaysWS, otherWS, err := FetchSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task)
	return c.Render(r.Context(), w)
}

func GetLoggedTime(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	workSessions, err := database.GetDatabase().GetWorkSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	for _, session := range workSessions {
		if session.OngoingToday {
			if (session.CreatedAt.Day() != time.Now().Day()) || (session.CreatedAt.Month() != time.Now().Month()) || (session.CreatedAt.Year() != time.Now().Year()) {
				if session.EndTime == nil {
					endTime := time.Date(session.StartTime.Year(), session.StartTime.Month(), session.StartTime.Day(), 23, 59, 59, 0, time.Local)
					session.EndTime = &endTime
					session.Duration += endTime.Sub(session.StartTime)
					session.LeftUnfinished = true
					session.Remaining = session.Remaining - session.Duration
					if session.Remaining < 0 {
						session.Remaining = 0
					}
				} else {
					session.Duration += session.EndTime.Sub(session.StartTime)
				}
				session.OngoingToday = false
				err = database.GetDatabase().UpdateWorkSession(&session)
				if err != nil {
					return err
				}
			}
		}
	}

	workSessions, err = database.GetDatabase().GetWorkSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	todaysWS := []database.WorkSession{}
	otherWS := []database.WorkSession{}
	for _, session := range workSessions {
		if session.OngoingToday {
			todaysWS = append(todaysWS, session)
		} else {
			otherWS = append(otherWS, session)
		}
	}

	otherWS = sortWorkSessionsByDate(otherWS)
	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task)

	return c.Render(r.Context(), w)

}

func GetChangeDuration(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}
	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}
	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	workSession, err := database.GetDatabase().GetWorkSessionByID(uint(sessionID))
	if err != nil {
		return err
	}

	c := ChangeDurationDialog(*workSession, uint(sprintID), uint(taskID), *task)

	return c.Render(r.Context(), w)

}

func PostChangeDuration(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}
	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	type RequestData struct {
		Duration string `schema:"duration,required"`
	}

	var req RequestData
	err = decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		w.WriteHeader(http.StatusSeeOther)
		c := common.ValidationError("Use XhXXm format for duration.")
		return c.Render(r.Context(), w)
	}

	if duration.Hours() > 24 {
		w.WriteHeader(http.StatusSeeOther)
		c := common.ValidationError("Duration cannot be more than 24 hours.")
		return c.Render(r.Context(), w)
	}

	if duration < 0 {
		w.WriteHeader(http.StatusSeeOther)
		c := common.ValidationError("Duration cannot be negative.")
		return c.Render(r.Context(), w)
	}

	session, err := database.GetDatabase().GetWorkSessionByID(uint(sessionID))
	if err != nil {
		return err
	}

	workSessions, err := database.GetDatabase().GetWorkSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	workSessions = sortWorkSessionsByDate(workSessions)

	var lastSession *database.WorkSession
	for _, ws := range workSessions {
		if ws.ID == session.ID {
			break
		}
		lastSession = &ws
	}

	remaining := time.Duration(0)

	if lastSession == nil {
		remaining = task.TimeEstimate - duration
		if remaining < 0 {
			remaining = 0
		}
	} else {
		remaining = lastSession.Remaining - duration
		if remaining < 0 {
			remaining = 0
		}
	}

	session.Duration = duration
	session.Remaining = remaining
	err = database.GetDatabase().UpdateWorkSession(session)
	if err != nil {
		return err
	}

	workSessions, err = database.GetDatabase().GetWorkSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	workSessions = sortWorkSessionsByDate(workSessions)

	todaysWS := []database.WorkSession{}
	otherWS := []database.WorkSession{}
	for _, session := range workSessions {
		if session.OngoingToday {
			todaysWS = append(todaysWS, session)
		} else {
			otherWS = append(otherWS, session)
		}
	}

	otherWS = sortWorkSessionsByDate(otherWS)

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task)

	return c.Render(r.Context(), w)

}


func GetChangeRemaining(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}
	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}
	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	workSession, err := database.GetDatabase().GetWorkSessionByID(uint(sessionID))
	if err != nil {
		return err
	}

	c := ChangeRemainingDialog(*workSession, uint(sprintID), uint(taskID), *task)

	return c.Render(r.Context(), w)

}

func PostChangeRemaining(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}
	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	type RequestData struct {
		Remaining string `schema:"remaining,required"`
	}

	var req RequestData
	err = decoder.Decode(&req, r.PostForm)
	if err != nil {
		return err
	}

	remaining, err := time.ParseDuration(req.Remaining)
	if err != nil {
		w.WriteHeader(http.StatusSeeOther)
		c := common.ValidationError("Use XhXXm format for duration.")
		return c.Render(r.Context(), w)
	}

	if remaining < 0 {
		w.WriteHeader(http.StatusSeeOther)
		c := common.ValidationError("Remaining time cannot be negative.")
		return c.Render(r.Context(), w)
	}

	session, err := database.GetDatabase().GetWorkSessionByID(uint(sessionID))
	if err != nil {
		return err
	}

	session.Remaining = remaining
	err = database.GetDatabase().UpdateWorkSession(session)
	if err != nil {
		return err
	}

	workSessions, err := database.GetDatabase().GetWorkSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	workSessions = sortWorkSessionsByDate(workSessions)

	todaysWS := []database.WorkSession{}
	otherWS := []database.WorkSession{}
	for _, session := range workSessions {
		if session.OngoingToday {
			todaysWS = append(todaysWS, session)
		} else {
			otherWS = append(otherWS, session)
		}
	}

	otherWS = sortWorkSessionsByDate(otherWS)

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task)

	return c.Render(r.Context(), w)

}

func GetUnfinishedSessionDialog(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	workSessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}

	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	workSession, err := database.GetDatabase().GetWorkSessionByID(uint(workSessionID))
	if err != nil {
		return err
	}

	workSession.LeftUnfinished = false
	err = database.GetDatabase().UpdateWorkSession(workSession)
	if err != nil {
		return err
	}

	c := UnfinishedSessionDialog(uint(taskID), uint(sprintID))
	return c.Render(r.Context(), w)
}

func FetchSessionsForTask(taskID uint) ([]database.WorkSession, []database.WorkSession, error) {
	workSessions, err := database.GetDatabase().GetWorkSessionsForTask(taskID)
	if err != nil {
		return nil, nil, err
	}

	todaysWS := []database.WorkSession{}
	otherWS := []database.WorkSession{}
	for _, session := range workSessions {
		if session.OngoingToday {
			todaysWS = append(todaysWS, session)
		} else {
			otherWS = append(otherWS, session)
		}
	}

	otherWS = sortWorkSessionsByDate(otherWS)

	return todaysWS, otherWS, nil
}
package task

import (
	"fmt"
	"net/http"
	"reflect"
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
		err = database.GetDatabase().CreateWorkSessionToday(params.UserID, uint(taskID))
		if err != nil {
			return err
		}
		w.Header().Set("HX-Redirect", fmt.Sprint("/sprintbacklog/", sprintID))
		w.WriteHeader(http.StatusSeeOther)

		return nil
	}

	err = database.GetDatabase().CreateWorkSessionToday(params.UserID, uint(taskID))
	if err != nil {
		return err
	}

	todaysWS, otherWS, err := FetchSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task, params.UserID)
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

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task, params.UserID)
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

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task, params.UserID)
	return c.Render(r.Context(), w)
}

func DeleteWorkSession(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
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

	err = database.GetDatabase().DeleteWorkSession(uint(workSessionID))
	if err != nil {
		return err
	}

	todaysWS, otherWS, err := FetchSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task, params.UserID)
	return c.Render(r.Context(), w)
}

func GetLoggedTime(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
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
	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}
	
	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task, params.UserID)

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

	// Fetch task again, so the data that was changed in trigger "AfterUpdate" will be updated
	task, err = database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task, params.UserID)

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

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task, params.UserID)

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

	c := UnfinishedSessionDialog(uint(sprintID), uint(taskID))
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

func GetStartPastWorkSession(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}

	sprintID, err := strconv.Atoi(chi.URLParam(r, "sprintID"))
	if err != nil {
		return err
	}

	c := createWorkSessionDialog(uint(sprintID), uint(taskID))
	return c.Render(r.Context(), w)
}

var timeConverter = func(value string) reflect.Value {
	layout := "2006-01-02"

	if v, err := time.Parse(layout, value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{}
}

type WorkSessionFormData struct {
	StartDate time.Time `schema:"start_date,required"`
	Duration  *float32  `schema:"duration,required"`
	Remaining *float32  `schema:"remaining,required"`
}

func PostStartPastWorkSession(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {

	var workSessionFormData WorkSessionFormData
	decoder.RegisterConverter(time.Time{}, timeConverter)
	if err := decoder.Decode(&workSessionFormData, r.PostForm); err != nil {
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

	task, err := database.GetDatabase().GetTaskByID(uint(taskID))
	if err != nil {
		return err
	}

	//check if session for this date already exists
	workSessions, err := database.GetDatabase().GetWorkSessionsForTask(uint(taskID))
	if err != nil {
		return err
	}

	for _, session := range workSessions {
		if session.StartTime.Day() == workSessionFormData.StartDate.Day() && session.StartTime.Month() == workSessionFormData.StartDate.Month() && session.StartTime.Year() == workSessionFormData.StartDate.Year() {
			w.WriteHeader(http.StatusSeeOther)
			c := common.ValidationError("Log for this date already exists. Go back and click to edit it.")
			return c.Render(r.Context(), w)
		}
	}

	//check if start date is before today7s dayso before the minight of today
	if workSessionFormData.StartDate.After(time.Now().Truncate(24 * time.Hour)) {
		w.WriteHeader(http.StatusSeeOther)
		c := common.ValidationError("Start date of a past log cannot be in the future.")
		return c.Render(r.Context(), w)
	}

	if workSessionFormData.StartDate.Equal(time.Now().Truncate(24 * time.Hour)) {
		w.WriteHeader(http.StatusSeeOther)
		c := common.ValidationError("Start date of a past log cannot be today. If you wish to start a log today, go back and click start.")
		return c.Render(r.Context(), w)
	}

	//check if start date is within the sprint
	sprint, err := database.GetDatabase().GetSprintByID(uint(sprintID))
	if err != nil {
		return err
	}

	if workSessionFormData.StartDate.Before(sprint.StartDate) || workSessionFormData.StartDate.After(sprint.EndDate) {
		w.WriteHeader(http.StatusSeeOther)
		c := common.ValidationError("Start date should be within the sprint.")
		return c.Render(r.Context(), w)
	}

	end := workSessionFormData.StartDate.Add(time.Duration(*workSessionFormData.Duration * float32(time.Hour)))

	// create new session
	session := database.WorkSession{
		StartTime:      workSessionFormData.StartDate,
		EndTime:        &end,
		Duration:       time.Duration(*workSessionFormData.Duration * float32(time.Hour)),
		Remaining:      time.Duration(*workSessionFormData.Remaining * float32(time.Hour)),
		TaskID:         uint(taskID),
		UserID:         params.UserID,
		OngoingToday:   false,
		LeftUnfinished: false,
	}

	//add session to database
	err = database.GetDatabase().CreateWorkSession(&session)
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

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID), uint(sprintID), *task, params.UserID)

	return c.Render(r.Context(), w)
}

func GetUserFirstAndLastNameFromID(userID uint) string {
	user, _ := database.GetDatabase().GetUserByID(userID)
	return user.FirstName + " " + user.LastName
}
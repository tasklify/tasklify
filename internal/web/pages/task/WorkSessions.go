package task

import (
	"net/http"
	"sort"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"
	"time"

	"github.com/go-chi/chi/v5"
)

func StartWorkSession(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
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
	}

	if err != nil {
		return err
	}

	err = database.GetDatabase().StartWorkSession(params.UserID, uint(taskID))
	if err != nil {
		return err
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(task.UserStoryID)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", "/sprintbacklog/"+strconv.Itoa(int(*userStory.SprintID)))

	w.WriteHeader(http.StatusSeeOther)
	return nil
}

func ResumeWorkSession(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	workSessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}

	workSession, err := database.GetDatabase().GetWorkSessionByID(uint(workSessionID))
	if err != nil {
		return err
	}

	err = database.GetDatabase().ResumeWorkSession(uint(workSessionID))
	if err != nil {
		return err
	}

	workSessions, err := database.GetDatabase().GetWorkSessionsForTask(workSession.TaskID)
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

	c := LoggedTimeDialog(todaysWS, otherWS, workSession.TaskID)
	return c.Render(r.Context(), w)
}

func StopWorkSession(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	workSessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}

	workSession, err := database.GetDatabase().GetWorkSessionByID(uint(workSessionID))
	if err != nil {
		return err
	}

	err = database.GetDatabase().EndWorkSession(uint(workSessionID))
	if err != nil {
		return err
	}

	workSessions, err := database.GetDatabase().GetWorkSessionsForTask(workSession.TaskID)
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

	c := LoggedTimeDialog(todaysWS, otherWS, workSession.TaskID)
	return c.Render(r.Context(), w)
}

func GetLoggedTime(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		return err
	}

	workSessions, err := database.GetDatabase().GetWorkSessionsForTask(uint(taskID))
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

	c := LoggedTimeDialog(todaysWS, otherWS, uint(taskID))

	return c.Render(r.Context(), w)

}

func sortWorkSessionsByDate(workSessions []database.WorkSession) []database.WorkSession {
	sort.Slice(workSessions, func(i, j int) bool {
		return workSessions[i].StartTime.Before(workSessions[j].StartTime)
	})

	return workSessions
}

func GetChangeDuration(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
	if err != nil {
		return err
	}

	workSession, err := database.GetDatabase().GetWorkSessionByID(uint(sessionID))
	if err != nil {
		return err
	}

	c := ChangeDurationDialog(*workSession)

	return c.Render(r.Context(), w)

}

func PostChangeDuration(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "workSessionID"))
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
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Use XhXXm format for duration.")
		return c.Render(r.Context(), w)
	}

	err = database.GetDatabase().ChangeDuration(uint(sessionID), duration)
	if err != nil {
		return err
	}

	workSession, err := database.GetDatabase().GetWorkSessionByID(uint(sessionID))
	if err != nil {
		return err
	}

	c := DurationDialog(*workSession)

	return c.Render(r.Context(), w)

}

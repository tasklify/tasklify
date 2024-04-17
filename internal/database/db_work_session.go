package database

import (
	"time"

	"gorm.io/gorm"
)

type WorkSession struct {
	gorm.Model
	StartTime      time.Time
	EndTime        *time.Time
	Duration       time.Duration
	Remaining      time.Duration
	TaskID         uint
	UserID         uint
	OngoingToday   bool
	LeftUnfinished bool
}

func (ws *WorkSession) AfterUpdate(tx *gorm.DB) (err error) {
	var workSessions = []WorkSession{}
	err = tx.Order("start_time asc").Find(&workSessions, "task_id = ?", ws.TaskID).Error
	if err != nil {
		return err
	}

	todaysWS := []WorkSession{}
	otherWS := []WorkSession{}
	for _, session := range workSessions {
		if session.OngoingToday {
			todaysWS = append(todaysWS, session)
		} else {
			otherWS = append(otherWS, session)
		}
	}

	task, err := databaseClient.GetTaskByID(ws.TaskID)
	if err != nil {
		return err
	}

	if (len(todaysWS) != 0 && todaysWS[0].Remaining == 0) || (len(todaysWS) == 0 && len(otherWS) != 0 && otherWS[len(otherWS)-1].Remaining == 0) {
		task.Status = &StatusDone
	} else {
		task.Status = &StatusInProgress
	}

	if err = databaseClient.UpdateTask(task); err != nil {
		return err
	}

	return nil
}

func (ws *WorkSession) BeforeDelete(tx *gorm.DB) (err error) {
	// Only data we provided on delete is available in ws, so get the whole object first
	err = tx.First(&ws, ws.ID).Error
	if err != nil {
		return err
	}

	// Fetch all work sessions but the one to be deleted
	var workSessions = []WorkSession{}
	err = tx.Order("start_time asc").Where("id != ?", ws.ID).Find(&workSessions, "task_id = ?", ws.TaskID).Error
	if err != nil {
		return err
	}

	todaysWS := []WorkSession{}
	otherWS := []WorkSession{}
	for _, session := range workSessions {
		if session.OngoingToday {
			todaysWS = append(todaysWS, session)
		} else {
			otherWS = append(otherWS, session)
		}
	}

	task, err := databaseClient.GetTaskByID(ws.TaskID)
	if err != nil {
		return err
	}

	if (len(todaysWS) != 0 && todaysWS[0].Remaining == 0) || (len(todaysWS) == 0 && len(otherWS) != 0 && otherWS[len(otherWS)-1].Remaining == 0) {
		task.Status = &StatusDone
	} else if len(todaysWS) == 0 && len(otherWS) == 0 {
		task.Status = &StatusTodo
	} else {
		task.Status = &StatusInProgress
	}

	if err = databaseClient.UpdateTask(task); err != nil {
		return err
	}

	return
}

func (db *database) GetWorkSessionByID(sessionID uint) (*WorkSession, error) {
	var session WorkSession
	err := db.First(&session, sessionID).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (db *database) CreateWorkSession(session *WorkSession) error {
	return db.Create(session).Error
}

func (db *database) CreateWorkSessionToday(userID, taskID uint) error {
	Remaining, err := db.GetRemainingTimeForTask(taskID)
	if err != nil {
		return err
	}
	session := WorkSession{
		StartTime:      time.Now(),
		TaskID:         taskID,
		UserID:         userID,
		OngoingToday:   true,
		LeftUnfinished: false,
		Remaining:      Remaining,
	}
	return db.Create(&session).Error
}

func (db *database) GetRemainingTimeForTask(taskID uint) (time.Duration, error) {
	var sessions []WorkSession
	err := db.Order("start_time asc").Find(&sessions, "task_id = ?", taskID).Error
	if err != nil {
		return 0, err
	}

	task, err := db.GetTaskByID(taskID)
	if err != nil {
		return 0, err
	}

	var defaultValue = task.TimeEstimate
	for _, session := range sessions {
		if !session.OngoingToday {
			defaultValue = session.Remaining
		}
	}

	return defaultValue, nil
}

func (db *database) UpdateWorkSession(session *WorkSession) error {
	return db.Save(session).Error
}

func (db *database) GetWorkSessionsForTask(taskID uint) ([]WorkSession, error) {
	var sessions []WorkSession
	err := db.Find(&sessions, "task_id = ?", taskID).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (db *database) DeleteWorkSession(sessionID uint) error {
	return db.Unscoped().Where("id = ?", sessionID).Delete(&WorkSession{Model: gorm.Model{ID: sessionID}}).Error
}

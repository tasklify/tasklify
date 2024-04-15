package database

import (
	"time"

	"gorm.io/gorm"
)

type WorkSession struct {
	gorm.Model
	StartTime    time.Time
	EndTime      *time.Time
	Duration     time.Duration
	Remaining    time.Duration
	TaskID       uint
	UserID       uint
	OngoingToday bool
	LeftUnfinished bool
}

func (db *database) GetWorkSessionByID(sessionID uint) (*WorkSession, error) {
	var session WorkSession
	err := db.First(&session, sessionID).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (db *database) CreateWorkSession(userID, taskID uint) error {
	Remaining, err := db.GetRemainingTimeForTask(taskID)
	if err != nil {
		return err
	}
	session := WorkSession{
		StartTime:    time.Now(),
		TaskID:       taskID,
		UserID:       userID,
		OngoingToday: true,
		LeftUnfinished: false,
		Remaining: Remaining,
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

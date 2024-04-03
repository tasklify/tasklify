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
	session := WorkSession{
		StartTime:    time.Now(),
		TaskID:       taskID,
		UserID:       userID,
		OngoingToday: true,
		LeftUnfinished: false,
	}
	return db.Create(&session).Error
}

func (db *database) UpdateWorkSession(session *WorkSession) error {
	return db.Save(session).Error
}

func (db *database) GetTotalTimeSpentOnTask(taskID uint) (time.Duration, error) {
	var sessions []WorkSession
	if err := db.Find(&sessions, "task_id = ?", taskID).Error; err != nil {
		return 0, err
	}

	var totalTime time.Duration
	for _, session := range sessions {
		totalTime += session.Duration
	}

	return totalTime, nil
}

func (db *database) GetWorkSessionsForTask(taskID uint) ([]WorkSession, error) {
	var sessions []WorkSession
	err := db.Find(&sessions, "task_id = ?", taskID).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

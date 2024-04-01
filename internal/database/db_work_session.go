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
}

func (db *database) GetWorkSessionByID(sessionID uint) (*WorkSession, error) {
	var session WorkSession
	err := db.First(&session, sessionID).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (db *database) StartWorkSession(userID, taskID uint) error {
	session := WorkSession{
		StartTime:    time.Now(),
		TaskID:       taskID,
		UserID:       userID,
		OngoingToday: true,
	}
	return db.Create(&session).Error
}

func (db *database) ResumeWorkSession(sessionID uint) error {
	var session WorkSession
	if err := db.First(&session, sessionID).Error; err != nil {
		return err
	}

	session.StartTime = time.Now()
	session.EndTime = nil

	return db.Save(&session).Error
}

func (db *database) EndWorkSession(sessionID uint) error {
	var session WorkSession
	if err := db.First(&session, sessionID).Error; err != nil {
		return err
	}

	endTime := time.Now()
	session.EndTime = &endTime
	session.Duration += endTime.Sub(session.StartTime)

	return db.Save(&session).Error
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

func (db *database) CloseOpenSessionsAtMidnight() error {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	var sessions []WorkSession
	err := db.Model(&WorkSession{}).Where("end_time IS NULL AND start_time < ?", midnight).Find(&sessions).Error
	if err != nil {
		return err
	}

	for _, session := range sessions {
		endTime := midnight
		session.EndTime = &endTime
		session.Duration += endTime.Sub(session.StartTime)
		session.OngoingToday = false
	}

	return nil
}

func (db *database) ChangeDuration(sessionID uint, duration time.Duration) error {
	var session WorkSession
	if err := db.First(&session, sessionID).Error; err != nil {
		return err
	}

	session.Duration = duration

	return db.Save(&session).Error
}

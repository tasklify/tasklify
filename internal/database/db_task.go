package database

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title          *string
	Description    *string `gorm:"type:TEXT"`
	TimeEstimate   time.Duration
	UserAccepted   *bool
	Status         *Status
	ProjectID      uint            // 1:n (ProjectHasUser:Task)
	UserID         *uint           // 1:n (ProjectHasUser:Task)
	ProjectHasUser *ProjectHasUser `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:Task)
	UserStoryID    uint            // 1:n (UserStory:Task)
	WorkSessions   []WorkSession   // 1:n (Task:WorkSession)
}

func (db *database) GetTasksByUserStory(userStoryID uint) ([]Task, error) {
	var tasks []Task

	err := db.Find(&tasks, "tasks.user_story_id = ?", userStoryID).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (db *database) CreateTask(task *Task) error {
	return db.Create(task).Error
}

func (db *database) UpdateTask(task *Task) error {
	return db.Save(task).Error
}

func (db *database) GetTaskByID(id uint) (*Task, error) {
	var task = &Task{}
	err := db.First(task, id).Error
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (db *database) DeleteTask(taskID uint) error {
	return db.Where("id = ?", taskID).Delete(&Task{}).Error
}

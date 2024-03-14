package database

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title          *string
	Description    *string `gorm:"type:TEXT"`
	TimeEstimate   *float32
	UserAccepted   *bool
	Status         *Status
	ProjectID      uint            // 1:n (ProjectHasUser:Task)
	UserID         *uint            // 1:n (ProjectHasUser:Task)
	ProjectHasUser *ProjectHasUser `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:Task)
	UserStoryID    uint            // 1:n (UserStory:Task)
}

func (db *database) GetTasksByUserStory(userStoryID uint) ([]Task, error) {
	var tasks []Task

	err := db.Find(&tasks, "tasks.user_story_id = ?", userStoryID).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

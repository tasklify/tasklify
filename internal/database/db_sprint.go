package database

import (
	"time"

	"gorm.io/gorm"
)

type Sprint struct {
	gorm.Model
	Title       string `gorm:"unique"`
	StartDate   time.Time
	EndDate     time.Time
	Velocity    *float32
	ProjectID   uint        // 1:n (Project:Sprint)
	UserStories []UserStory // 1:n (Sprint:UserStory)
}

func (db *database) CreateSprint(sprint *Sprint) error {
	return db.Create(sprint).Error
}

func (db *database) GetSprintByProject(projectID uint) ([]Sprint, error) {
	var sprints []Sprint

	err := db.Find(&sprints, "sprints.project_id = ?", "1").Error
	if err != nil {
		return nil, err
	}

	return sprints, nil
}

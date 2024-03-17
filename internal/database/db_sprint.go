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
	UserStories []UserStory `gorm:"foreignKey:SprintID"` // 1:n (Sprint:UserStory)
}

func (db *database) CreateSprint(sprint *Sprint) error {
	return db.Create(sprint).Error
}

func (db *database) GetSprintByProject(projectID uint) ([]Sprint, error) {
	var sprints []Sprint

	err := db.Preload("UserStories.AcceptanceTests").Preload("UserStories").Find(&sprints, "sprints.project_id = ?", projectID).Error
	if err != nil {
		return nil, err
	}

	return sprints, nil
}

func (sprint *Sprint) DetermineStatus() Status {
	now := time.Now()

	//todo is this correct be careful of edge conditions
	if now.Before(sprint.StartDate) {
		return StatusTodo
	} else if now.After(sprint.StartDate) && now.Before(sprint.EndDate) {
		return StatusInProgress
	} else if now.After(sprint.EndDate) {
		return StatusDone
	} else {
		return StatusTodo
	}
}

func (db *database) GetSprintByID(id uint) (*Sprint) {
	var sprint = &Sprint{}
	err := db.First(sprint, id).Error
	if err != nil {
		return nil
	}

	return sprint
}
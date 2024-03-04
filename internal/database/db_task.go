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
	Status         *StatusEnum     `gorm:"type:status_enum"`
	ProjectHasUser *ProjectHasUser `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:Task)
	UserStoryID    uint            // 1:n (UserStory:Task)
}

type StatusEnum string

const (
	StatusTodo       StatusEnum = "todo"
	StatusInProgress StatusEnum = "in progress"
	StatusDone       StatusEnum = "done"
)

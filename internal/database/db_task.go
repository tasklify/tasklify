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
	UserID         uint            // 1:n (ProjectHasUser:Task)
	ProjectHasUser *ProjectHasUser `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:Task)
	UserStoryID    uint            // 1:n (UserStory:Task)
}

package database

import (
	"time"

	"gorm.io/gorm"
)

// This is a JoinTable for m:n relation between Project and User
type ProjectHasUser struct {
	ProjectID   uint `gorm:"primaryKey;autoIncrement:false"`
	UserID      uint `gorm:"primaryKey;autoIncrement:false"`
	ProjectRole ProjectRole
	Tasks       []Task      `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:Task) -> Every user that works on a project can work on multiple tasks
	UserStories []UserStory `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:UserStory) -> Every user that works on a project can be owner of multiple user stories
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

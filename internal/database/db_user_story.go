package database

import (
	"gorm.io/gorm"
)

type UserStory struct {
	gorm.Model
	Title            string  `gorm:"unique"`
	Description      *string `gorm:"type:TEXT"`
	Priority         Priority
	BusinessValue    int
	Realized         *bool
	RejectionComment *string `gorm:"type:TEXT"`
	WorkflowStepID   *uint   // 1:1 (WorkflowStep:UserStory)
	WorkflowStep     WorkflowStep
	SprintID         *uint           // 1:n (Sprint:UserStory)
	ProjectID        uint            // 1:n (Project:UserStory)
	Tasks            []Task          // 1:n (UserStory:Task)
	UserID           uint            // 1:n (ProjectHasUser:UserStory)
	ProjectHasUser   *ProjectHasUser `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:UserStory)
}

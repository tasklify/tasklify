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
	StoryNumber      int             `gorm:"not null"`
}

func (db *database) CreateUserStory(userStory *UserStory) error {
	var count int64
	db.Model(&UserStory{}).Where("project_id = ?", userStory.ProjectID).Count(&count)

	// Set the StoryNumber to be the next number in sequence for the project
	userStory.StoryNumber = int(count) + 1

	return db.Create(userStory).Error
}

func (db *database) GetUserStoryByID(id uint) (*UserStory, error) {
	var userStory = &UserStory{}
	err := db.First(userStory, id).Error
	if err != nil {
		return nil, err
	}

	return userStory, nil
}

func (db *database) GetUserStoryByProject(projectID uint) ([]UserStory, error) {
	var userStories []UserStory
	err := db.Where("project_id = ?", projectID).Find(&userStories).Error
	if err != nil {
		return []UserStory{}, err
	}

	return userStories, nil
}

func (db *database) GetUserStoryByUser(userID uint) ([]UserStory, error) {
	var userStories []UserStory
	err := db.Where("user_id = ?", userID).Find(&userStories).Error
	if err != nil {
		return []UserStory{}, err
	}

	return userStories, nil
}

func (db *database) UserStoryWithTitleExists(title string) bool {
	var count int64
	db.Model(&UserStory{}).Where("title = ?", title).Count(&count)
	return count > 0
}

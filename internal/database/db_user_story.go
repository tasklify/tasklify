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
	UserID           *uint            // 1:n (ProjectHasUser:UserStory)
	ProjectHasUser   *ProjectHasUser `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:UserStory)
}

func (db *database) CreateUserStory(userStory *UserStory) error {
	return db.Create(userStory).Error
}

func (db *database) GetUserStoriesByProject(projectID uint) ([]UserStory, error) {
	var userStories []UserStory

	err := db.Find(&userStories, "user_stories.project_id = ?", projectID).Error
	if err != nil {
		return nil, err
	}

	return userStories, nil
}

func (db *database) GetUserStoriesBySprint(sprintID uint) ([]UserStory, error) {
	var userStories []UserStory

	err := db.Find(&userStories, "user_stories.sprint_id = ?", sprintID).Error
	if err != nil {
		return nil, err
	}

	return userStories, nil
}

func (db *database) GetUserStoryByID(id uint) (*UserStory, error) {
	var userStory = &UserStory{}
	err := db.First(userStory, id).Error
	if err != nil {
		return nil, err
	}

	return userStory, nil
}

func (db *database) UserStoryWithTitleExists(title string) bool {
	var count int64
	db.Model(&UserStory{}).Where("title = ?", title).Count(&count)
	return count > 0
}

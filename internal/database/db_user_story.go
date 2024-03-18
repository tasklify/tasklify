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
	StoryPoints      uint
	Realized         *bool
	RejectionComment *string `gorm:"type:TEXT"`
	WorkflowStepID   *uint   // 1:1 (WorkflowStep:UserStory)
	WorkflowStep     WorkflowStep
	SprintID         *uint           // 1:n (Sprint:UserStory)
	ProjectID        uint            // 1:n (Project:UserStory)
	Tasks            []Task          // 1:n (UserStory:Task)
	AcceptanceTests  []AcceptanceTest // 1:n (UserStory:AcceptanceTest)
	UserID           *uint           // 1:n (ProjectHasUser:UserStory)
	ProjectHasUser   *ProjectHasUser `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:UserStory)
}

func (db *database) CreateUserStory(userStory *UserStory) error {
	return db.Create(userStory).Error
}

func (db *database) UpdateUserStory(userStory *UserStory) error {
	return db.Save(userStory).Error
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
	db.Model(&UserStory{}).Where("LOWER(title) = LOWER(?)", title).Count(&count)
	return count > 0
}

func (db *database) AddUserStoryToSprint(sprintID uint, userStoryIDs []uint) (*Sprint, error) {
	// Update user stories with the sprint ID
	if err := db.Model(&UserStory{}).Where("id IN (?)", userStoryIDs).Update("sprint_id", sprintID).Error; err != nil {
		return nil, err
	}

	// Retrieve the sprint
	var sprint Sprint
	if err := db.Preload("UserStories").First(&sprint, sprintID).Error; err != nil {
		return nil, err
	}

	return &sprint, nil
}

func (userStory *UserStory) AllAcceptanceTestsRealized() bool {
	for _, acceptanceTest := range userStory.AcceptanceTests {
		if *acceptanceTest.Realized == false {
			return false
		}
	}
	return true
}

func (userStory *UserStory) AllTasksRealized() bool {
	for _, task := range userStory.Tasks {
		if *task.Status != StatusDone {
			return false
		}
	}
	return true
}
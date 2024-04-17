package database

import (
	"gorm.io/gorm"
)

type AcceptanceTest struct {
	gorm.Model
	Description    *string `gorm:"type:TEXT"`
	Realized         *bool
	UserStoryID      uint
}

func (db *database) CreateAcceptanceTest(acceptanceTest *AcceptanceTest) error {
	return db.Create(&acceptanceTest).Error
}

func (db *database) GetAcceptanceTestsByUserStory(userStoryID uint) ([]AcceptanceTest, error) {
	var acceptanceTests []AcceptanceTest

	err := db.Find(&acceptanceTests, "acceptance_tests.user_story_id = ?", userStoryID).Error
	if err != nil {
		return nil, err
	}

	return acceptanceTests, nil
}

func (db *database) DeleteAcceptanceTest(acceptanceTest *AcceptanceTest) error {
	return db.Unscoped().Delete(&acceptanceTest).Error
}
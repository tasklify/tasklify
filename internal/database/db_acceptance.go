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

func (db *database) GetAcceptanceTestByID(id uint) (*AcceptanceTest, error) {
	var acceptanceTest AcceptanceTest

	err := db.First(&acceptanceTest, id).Error
	if err != nil {
		return nil, err
	}

	return &acceptanceTest, nil
}

func (db *database) UpdateAcceptanceTest(acceptanceTest *AcceptanceTest) error {
	return db.Save(&acceptanceTest).Error
}
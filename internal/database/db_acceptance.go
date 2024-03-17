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
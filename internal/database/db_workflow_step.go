package database

import (
	"gorm.io/gorm"
)

type WorkflowStep struct {
	gorm.Model
	Key   string `gorm:"unique"`
	Title string
}

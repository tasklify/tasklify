package database

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Title       string      `gorm:"unique"`
	Users       []User      `gorm:"many2many:project_has_users;"` // m:n (Project:User)
	Sprints     []Sprint    // 1:n (Project:Sprint)
	UserStories []UserStory // 1:n (Project:UserStory)
}

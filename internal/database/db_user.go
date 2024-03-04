package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	Password     string
	FirstName    string
	LastName     string
	Email        string `gorm:"unique"`
	LastLogin    *time.Time
	SystemRoleID uint // 1:1 (User:SystemRole)
	SystemRole   SystemRole
	Projects     []Project `gorm:"many2many:project_has_users;"` // m:n (User:Project)
}

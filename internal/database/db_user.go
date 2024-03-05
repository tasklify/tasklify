package database

import (
	"crypto/sha256"
	"fmt"
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

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// Store hashed password
	if u.Password != "" {
		h := sha256.Sum256([]byte(u.Password))
		u.Password = fmt.Sprintf("%x", h[:])
	}

	return nil
}

func (d *database) CreateUser(user *User) error {
	return d.Create(user).Error
}

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

func (db *database) GetUser(username string) (*User, error) {
	var user = &User{Username: username}
	err := db.First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *database) CreateUser(user *User) error {
	return db.Create(user).Error
}

func (db *database) UpdateUser(user *User) error {
	return db.Save(user).Error
}

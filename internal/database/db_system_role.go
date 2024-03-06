package database

import (
	"gorm.io/gorm"
)

type SystemRole struct {
	gorm.Model
	Key   string `gorm:"unique"`
	Title string
}

func (db *database) GetSystemRole(systemRoleName string) (*SystemRole, error) {
	var systemRole = &SystemRole{
		Key: systemRoleName,
	}

	err := db.First(systemRole).Error
	if err != nil {
		return nil, err
	}

	return systemRole, nil
}

func (db *database) CreateSystemRole(role *SystemRole) error {
	return db.Create(role).Error
}

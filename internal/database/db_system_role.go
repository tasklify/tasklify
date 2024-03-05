package database

import (
	"gorm.io/gorm"
)

type SystemRole struct {
	gorm.Model
	Key   string `gorm:"unique"`
	Title string
}

func (d *database) CreateSystemRole(role *SystemRole) error {
	return d.Create(role).Error
}

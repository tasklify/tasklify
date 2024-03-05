package database

import (
	"gorm.io/gorm"
)

type ProjectRole struct {
	gorm.Model
	Key   string `gorm:"unique"`
	Title string
}

func (d *database) CreateProjectRole(role *ProjectRole) error {
	return d.Create(role).Error
}

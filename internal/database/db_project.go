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

func (db *database) GetProjectByID(id uint) (*Project, error) {
	var project = &Project{}
	err := db.First(project, id).Error
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (db *database) CreateProject(project *Project) (uint, error) {
	err := db.Create(&project).Error
	return project.ID, err
}

func (db *database) ProjectWithTitleExists(title string) bool {
	var count int64
	db.Model(&Project{}).Where("title = ?", title).Count(&count)
	return count > 0
}

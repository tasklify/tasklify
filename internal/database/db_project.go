package database

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Title       string      `gorm:"unique"`
	Description string      `gorm:"type:TEXT"`
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
	db.Model(&Project{}).Where("LOWER(title) = LOWER(?)", title).Count(&count)
	return count > 0
}

func (db *database) GetUserProjects(userID uint) ([]Project, error) {
	user, err := db.GetUserByID(userID)
	if err != nil {
		return []Project{}, err
	}

	if user.SystemRole == SystemRoleAdmin {
		var projects []Project
		if err := db.Find(&projects).Error; err != nil {
			return []Project{}, err
		}
		return projects, nil
	} else {
		if err := db.Preload("Projects").First(&user, userID).Error; err != nil {
			return []Project{}, err
		}
		return user.Projects, nil
	}
}

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

func (db *database) ProjectWithTitleExists(title string, excludedProjectID *uint) bool {
	var count int64
	if excludedProjectID != nil {
		db.Model(&Project{}).Where("LOWER(title) = LOWER(?) AND id != ?", title, excludedProjectID).Count(&count)
	} else {
		db.Model(&Project{}).Where("LOWER(title) = LOWER(?)", title).Count(&count)
	}
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
		var projects = []Project{}
		if err := db.
			Joins("JOIN project_has_users ON projects.id = project_has_users.project_id").
			Where("project_has_users.user_id = ? AND project_has_users.deleted_at IS NULL", userID).
			Find(&projects).Error; err != nil {
			return nil, err
		}
		return projects, nil
	}
}

func (db *database) UpdateProject(projectID uint, projectData Project) error {
	return db.Model(&Project{}).Where("id = ?", projectID).Updates(projectData).Error
}

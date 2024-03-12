package database

import (
	"time"

	"gorm.io/gorm"
)

// This is a JoinTable for m:n relation between Project and User
type ProjectHasUser struct {
	ProjectID   uint `gorm:"primaryKey;autoIncrement:false"`
	UserID      uint `gorm:"primaryKey;autoIncrement:false"`
	ProjectRole ProjectRole
	Tasks       []Task      `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:Task) -> Every user that works on a project can work on multiple tasks
	UserStories []UserStory `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:UserStory) -> Every user that works on a project can be owner of multiple user stories
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (db *database) AddUserToProject(projectID uint, userID uint, projectRole string) error {
	projectHasUser := ProjectHasUser{
		ProjectID:   projectID,
		UserID:      userID,
		ProjectRole: *ProjectRoles.Parse(projectRole),
	}

	if err := db.Create(&projectHasUser).Error; err != nil {
		return err
	}

	return nil
}

func (db *database) GetUsersOnProject(projectID uint) ([]User, error) {
	var project Project
	if err := db.Preload("Users").First(&project, projectID).Error; err != nil {
		return nil, err
	}

	for i := range project.Users {
		var projectHasUser ProjectHasUser
		if err := db.Where("project_id = ? AND user_id = ?", projectID, project.Users[i].ID).First(&projectHasUser).Error; err != nil {
			return nil, err
		}

		project.Users[i].ProjectRole = projectHasUser.ProjectRole
	}

	return project.Users, nil
}

func (db *database) GetUsersNotOnProject(projectID uint) ([]User, error) {
	var users []User
	subQuery := db.Model(&ProjectHasUser{}).Select("user_id").Where("user_id IS NOT NULL")
	if err := db.
		Not("id IN (?)", subQuery).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

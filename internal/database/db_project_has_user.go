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

// This struct is used so we can retrieve ProjectRole in addition to all User fields
type UserWithRole struct {
	User
	ProjectRoleStr string
}

func (db *database) AddUserToProject(projectID uint, userID uint, projectRole string) error {

	// If user has been previously removed from project (soft deleted record), we need to add it back by setting field deleted_at to nil
	var count int64
	db.Unscoped().Model(&ProjectHasUser{}).Where("project_id = ? AND user_id = ?", projectID, userID).Count(&count)
	if count == 1 {
		if err := db.Unscoped().
			Model(&ProjectHasUser{}).
			Where("project_id = ? AND user_id = ?", projectID, userID).
			Updates(map[string]interface{}{"deleted_at": nil, "project_role": *ProjectRoles.Parse(projectRole)}).Error; err != nil {
			return err
		}
		return nil
	}

	// If user has not been added to project before, we have to create new record
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
	var dbUsers []UserWithRole
	if err := db.Model(&User{}).
		Joins("LEFT JOIN project_has_users ON users.id = project_has_users.user_id").
		Select("users.*, project_has_users.project_role as project_role_str").
		Where("project_has_users.project_id = ? AND project_has_users.deleted_at IS NULL", projectID).
		Find(&dbUsers).Error; err != nil {
		return []User{}, err
	}

	// TODO: maybe in the future, we can find a beter way of getting project role for each user
	var users []User
	for _, u := range dbUsers {
		user := u.User
		parsedProjectRole := ProjectRoles.Parse(u.ProjectRoleStr)
		if parsedProjectRole != nil {
			user.ProjectRole = *parsedProjectRole
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *database) GetUsersWithRoleOnProject(projectID uint, projectRole ProjectRole) ([]User, error) {
	var users []User
	if err := db.Model(&User{}).
		Joins("LEFT JOIN project_has_users ON users.id = project_has_users.user_id").
		Select("users.*, project_has_users.project_role as project_role_str").
		Where("project_has_users.project_id = ? AND project_has_users.project_role = ? AND project_has_users.deleted_at IS NULL", projectID, projectRole.Val).
		Find(&users).Error; err != nil {
		return []User{}, err
	}

	for i := range users {
		users[i].ProjectRole = projectRole
	}

	return users, nil
}

func (db *database) GetUsersNotOnProject(projectID uint) ([]User, error) {
	var users []User

	if err := db.
		Select("users.*").
		Joins("LEFT JOIN project_has_users AS p ON users.id = p.user_id AND p.project_id = ?", projectID).
		Where("p.user_id IS NULL OR p.deleted_at IS NOT NULL").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (db *database) RemoveUserFromProject(projectID uint, userID uint) error {
	projectHasUser := ProjectHasUser{
		ProjectID: projectID,
		UserID:    userID,
	}

	if err := db.Delete(&projectHasUser).Error; err != nil {
		return err
	}

	return nil
}

func (db *database) GetProjectRole(userID uint, projectID uint) (ProjectRole, error) {
	var projectHasUser ProjectHasUser
	db.Where("user_id = ? AND project_id = ?", userID, projectID).First(&projectHasUser)

	return projectHasUser.ProjectRole, nil
}

func (db *database) GetProjectHasUserByProjectAndUser(userID uint, projectID uint) (*ProjectHasUser, error) {

	var projectHasUser = &ProjectHasUser{ProjectID: projectID, UserID: userID}
	err := db.First(projectHasUser).Error
	if err != nil {
		return nil, err
	}

	return projectHasUser, nil
}

func (db *database) UpsertUserOnProject(projectID uint, userID uint, projectRole string) error {
	var count int64

	// If user already exists on project, just update it with new data
	db.Model(&ProjectHasUser{}).Where("project_id = ? AND user_id = ?", projectID, userID).Count(&count)
	if count == 1 {
		err := db.Model(&ProjectHasUser{}).Where("project_id = ? AND user_id = ?", projectID, userID).Update("project_role", projectRole).Error
		if err != nil {
			return err
		}

	} else {
		err := db.AddUserToProject(projectID, userID, projectRole)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *database) RemoveUsersNotInList(projectID uint, userIDs []uint) error {
	return db.Where("project_id = ? AND user_id NOT IN ?", projectID, userIDs).Delete(&ProjectHasUser{}).Error
}

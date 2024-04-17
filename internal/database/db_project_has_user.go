package database

import (
	"time"

	"gorm.io/gorm"
)

// This is a JoinTable for m:n relation between Project and User
type ProjectHasUser struct {
	ProjectID uint `gorm:"primaryKey;autoIncrement:false"`
	UserID    uint `gorm:"primaryKey;autoIncrement:false"`
	// ProjectRole ProjectRole
	Tasks       []Task      `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:Task) -> Every user that works on a project can work on multiple tasks
	UserStories []UserStory `gorm:"foreignKey:ProjectID,UserID"` // 1:n (ProjectHasUser:UserStory) -> Every user that works on a project can be owner of multiple user stories
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (db *database) AddDeveloperToProject(projectID uint, userID uint) error {

	// If user has been previously removed from project (soft deleted record), we need to add it back by setting field deleted_at to nil
	var count int64
	db.Unscoped().Model(&ProjectHasUser{}).Where("project_id = ? AND user_id = ?", projectID, userID).Count(&count)
	if count == 1 {
		if err := db.Unscoped().
			Model(&ProjectHasUser{}).
			Where("project_id = ? AND user_id = ?", projectID, userID).
			Updates(map[string]interface{}{"deleted_at": nil}).Error; err != nil {
			return err
		}
		return nil
	}

	// If user has not been added to project before, we have to create new record
	projectHasUser := ProjectHasUser{
		ProjectID: projectID,
		UserID:    userID,
	}

	if err := db.Create(&projectHasUser).Error; err != nil {
		return err
	}

	return nil
}

func (db *database) GetUsersWithRoleOnProject(projectID uint, projectRole ProjectRole) ([]User, error) {
	var users []User
	if projectRole == ProjectRoleManager {
		if err := db.Model(&User{}).
			Joins("INNER JOIN projects ON users.id = projects.product_owner_id").
			Select("users.*").
			Where("projects.id = ?", projectID).
			Find(&users).Error; err != nil {
			return []User{}, err
		}
	} else if projectRole == ProjectRoleMaster {
		if err := db.Model(&User{}).
			Joins("INNER JOIN projects ON users.id = projects.scrum_master_id").
			Select("users.*").
			Where("projects.id = ?", projectID).
			Find(&users).Error; err != nil {
			return []User{}, err
		}
	} else if projectRole == ProjectRoleDeveloper {
		if err := db.Model(&User{}).
			Joins("LEFT JOIN project_has_users ON users.id = project_has_users.user_id").
			Select("users.*").
			Where("project_has_users.project_id = ? AND project_has_users.deleted_at IS NULL", projectID).
			Find(&users).Error; err != nil {
			return []User{}, err
		}
	}

	return users, nil
}

// Returns all users that can be added to project as developers (this includes scrum master, because he can have both of those roles)
func (db *database) GetUsersNotOnProject(projectID uint) ([]User, error) {
	var users []User

	if err := db.
		Select("users.*").
		Joins("LEFT JOIN project_has_users AS pu ON users.id = pu.user_id AND pu.project_id = ?", projectID).
		Joins("LEFT JOIN projects AS pr ON users.id = pr.product_owner_id AND pr.id = ?", projectID).
		Where("(pu.user_id IS NULL OR pu.deleted_at IS NOT NULL) AND pr.product_owner_id IS NULL").
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

func (db *database) GetProjectRoles(userID uint, projectID uint) ([]ProjectRole, error) {
	var count int64

	// Check if user is product owner, then return because in that case he can only have one role
	db.Model(&Project{}).Where("id = ? AND product_owner_id = ?", projectID, userID).Count(&count)
	if count == 1 {
		return []ProjectRole{ProjectRoleManager}, nil
	}

	var projectRoles []ProjectRole

	// Otherwise check if user is developer
	db.Model(&ProjectHasUser{}).Where("project_id = ? AND user_id = ?", projectID, userID).Count(&count)
	if count == 1 {
		projectRoles = append(projectRoles, ProjectRoleDeveloper)
	}

	// And check if user is also scrum master
	db.Model(&Project{}).Where("id = ? AND scrum_master_id = ?", projectID, userID).Count(&count)
	if count == 1 {
		projectRoles = append(projectRoles, ProjectRoleMaster)
	}

	return projectRoles, nil
}

func (db *database) GetProjectHasUserByProjectAndUser(userID uint, projectID uint) (*ProjectHasUser, error) {

	var projectHasUser = &ProjectHasUser{ProjectID: projectID, UserID: userID}
	err := db.First(projectHasUser).Error
	if err != nil {
		return nil, err
	}

	return projectHasUser, nil
}

// func (db *database) UpsertUserOnProject(projectID uint, userID uint, projectRole string) error {
// 	var count int64

// 	// If user already exists on project, just update it with new data
// 	db.Model(&ProjectHasUser{}).Where("project_id = ? AND user_id = ?", projectID, userID).Count(&count)
// 	if count == 1 {
// 		err := db.Model(&ProjectHasUser{}).Where("project_id = ? AND user_id = ?", projectID, userID).Update("project_role", projectRole).Error
// 		if err != nil {
// 			return err
// 		}

// 	} else {
// 		err := db.AddDeveloperToProject(projectID, userID, projectRole)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (db *database) RemoveUsersNotInList(projectID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		projectHasUsers := []ProjectHasUser{}
		if err := db.Where("project_id = ?", projectID).Preload("Tasks").Find(&projectHasUsers).Error; err != nil {
			return err
		}
		return db.Where("project_id = ?", projectID).Delete(&projectHasUsers).Error
	}

	projectHasUsers := []ProjectHasUser{}
	if err := db.Where("project_id = ? AND user_id NOT IN ?", projectID, userIDs).Preload("Tasks").Find(&projectHasUsers).Error; err != nil {
		return err
	}

	return db.Where("project_id = ? AND user_id NOT IN ?", projectID, userIDs).Delete(&projectHasUsers).Error
}

func (phu *ProjectHasUser) BeforeDelete(tx *gorm.DB) (err error) {
	if err := phu.UnclaimUserTasks(tx); err != nil {
		return err
	}

	return
}

func (phu *ProjectHasUser) UnclaimUserTasks(tx *gorm.DB) (err error) {
	for _, task := range phu.Tasks {
		// First, stop active time logging on this task
		activeWS := []WorkSession{}
		if err := tx.Where("user_id = ? AND end_time IS NULL AND task_id = ?", phu.UserID, task.ID).Find(&activeWS).Error; err != nil {
			return err
		}

		for _, ws := range activeWS {
			endTime := time.Now()
			ws.Remaining = ws.Remaining - endTime.Sub(ws.StartTime)
			if ws.Remaining < 0 {
				ws.Remaining = 0
			}
			ws.EndTime = &endTime
			ws.Duration += endTime.Sub(ws.StartTime)
			ws.LeftUnfinished = false

			if err = tx.Save(&ws).Error; err != nil {
				return err
			}
		}

		// Then, unclaim the task
		task.UserAccepted = new(bool)
		task.UserID = nil

		if err = tx.Save(&task).Error; err != nil {
			return err
		}
	}

	return
}

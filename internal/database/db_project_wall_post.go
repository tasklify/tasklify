package database

import (
	"gorm.io/gorm"
)

type ProjectWallPost struct {
	gorm.Model
	ProjectID uint   // 1:n (Project:ProjectWallPost)
	AuthorID  uint   // 1:1 (ProjectWallPost:User) relation Belongs To
	Author    User   `gorm:"foreignKey:AuthorID"` // 1:1 (ProjectWallPost:User)
	Body      string `gorm:"type:TEXT"`
}

func (db *database) GetProjectWallPosts(projectID uint) ([]ProjectWallPost, error) {
	var posts = []ProjectWallPost{}
	err := db.Preload("Author").Where("project_id = ?", projectID).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return []ProjectWallPost{}, err
	}

	return posts, nil
}

func (db *database) GetProjectWallPostByID(postID uint) (*ProjectWallPost, error) {
	var post = &ProjectWallPost{}
	err := db.First(post, postID).Error
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (db *database) AddProjectWallPost(post ProjectWallPost) error {
	return db.Create(&post).Error
}

func (db *database) EditProjectWallPost(projectID uint, postID uint, body string) error {
	return db.Model(&ProjectWallPost{}).Where("id = ? AND project_id = ?", postID, projectID).Update("body", body).Error
}

func (db *database) DeleteProjectWallPost(projectID uint, postID uint) error {
	return db.Unscoped().Where("id = ? AND project_id = ?", postID, projectID).Delete(&ProjectWallPost{}).Error
}

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
	err := db.Preload("Author").Where("project_id = ?", projectID).Order("created_at ASC").Find(&posts).Error
	if err != nil {
		return []ProjectWallPost{}, err
	}

	return posts, nil
}

func (db *database) AddProjectWallPost(post ProjectWallPost) error {
	return db.Create(&post).Error
}

func (db *database) EditProjectWallPost(postID uint, body string) error {
	return db.Model(&ProjectWallPost{}).Where("id = ?", postID).Update("body", body).Error
}

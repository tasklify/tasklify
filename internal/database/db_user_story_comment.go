package database

import "gorm.io/gorm"

type UserStoryComment struct {
	gorm.Model
	UserStoryID uint   // 1:n (UserStory:UserStoryComment)
	AuthorID    uint   // 1:1 (UserStoryComment:User) relation Belongs To
	Author      User   `gorm:"foreignKey:AuthorID"` // 1:1 (UserStoryComment:User)
	Body        string `gorm:"type:TEXT"`
}

func (db *database) GetUserStoryComments(userStoryID uint) ([]UserStoryComment, error) {
	var comments = []UserStoryComment{}
	err := db.Preload("Author").Where("user_story_id = ?", userStoryID).Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return []UserStoryComment{}, err
	}

	return comments, nil
}

func (db *database) GetUserStoryCommentByID(commentID uint) (*UserStoryComment, error) {
	var comment = &UserStoryComment{}
	err := db.First(comment, commentID).Error
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (db *database) AddUserStoryComment(comment UserStoryComment) error {
	return db.Create(&comment).Error
}

func (db *database) EditUserStoryComment(userStoryID uint, commentID uint, body string) error {
	return db.Model(&UserStoryComment{}).Where("id = ? AND user_story_id = ?", commentID, userStoryID).Update("body", body).Error
}

func (db *database) DeleteUserStoryComment(userStoryID uint, commentID uint) error {
	return db.Unscoped().Where("id = ? AND user_story_id = ?", commentID, userStoryID).Delete(&UserStoryComment{}).Error
}

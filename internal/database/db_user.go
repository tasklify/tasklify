package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"unique"`
	Password    string
	TotpURL     *string
	FirstName   string
	LastName    string
	Email       string `gorm:"unique"`
	LastLogin   *time.Time
	SystemRole  SystemRole  `gorm:"type:string"`
	Projects    []Project   `gorm:"many2many:project_has_users;"` // m:n (User:Project)
	ProjectRole ProjectRole `gorm:"type:string;-"`                // This field is ignored in database and all queries
}

// If you specify callerUserID it will get execluded from list
func (db *database) GetUsers() ([]User, error) {
	var users []User

	err := db.Find(&users).Error
	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func (db *database) GetFilteredUsers(userIDs []uint) ([]User, error) {
	var users = []User{}
	err := db.Where("id NOT IN ?", userIDs).Find(&users).Error
	if err != nil {
		return []User{}, err
	}

	fmt.Println(users)

	return users, nil
}

func (db *database) GetUserByUsername(username string) (*User, error) {
	var user = &User{}
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *database) GetUserByID(id uint) (*User, error) {
	var user = &User{}
	user.ID = id
	err := db.First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *database) DeleteUserByID(id uint) error {
	var user = &User{}
	user.ID = id
	return db.Delete(user).Error
}

func (db *database) CreateUser(user *User) error {
	return db.Create(user).Error
}

func (db *database) UpdateUser(user *User) error {
	return db.Save(user).Error
}

func (db *database) DeleteUser(user *User) error {
	return db.Delete(user).Error
}

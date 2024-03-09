package database

import (
	"fmt"
	"log"
	"sync"
	"tasklify/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	GetUser(username string) (*User, error)
	UpdateUser(user *User) error
	CreateUser(user *User) error
	CreateSprint(sprint *Sprint) error
	CreateUserStory(userStory *UserStory) error
	// GetSystemRole(systemRoleName string) (*SystemRole, error)
	RawDB() *gorm.DB
}

type database struct {
	*gorm.DB
}

var (
	onceDatabase sync.Once

	databaseClient *database
)

func GetDatabase(config ...*config.Config) Database {

	onceDatabase.Do(func() { // <-- atomic, does not allow repeating
		config := config[0]

		databaseClient = connectDatabase(config.Database)
		registerTables(databaseClient)
		populateDatabase(databaseClient)
	})

	return databaseClient
}

func connectDatabase(config config.Database) *database {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host, config.User, config.Password, config.DbName, config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected")

	return &database{DB: db}
}

func registerTables(db *database) {
	// Migrate the schema
	err := db.AutoMigrate(&User{}, &Project{}, &UserStory{}, &Task{}, &ProjectHasUser{}, &Sprint{}, &WorkflowStep{})
	if err != nil {
		log.Fatal("Schema migration error: ", err)
	}

	log.Println("Database tables registered")
}

func populateDatabase(db *database) {
	var count int64

	// Create Users
	if db.Model(&User{}).Count(&count); count == 0 {
		for _, user := range users {
			if err := db.CreateUser(&user); err != nil {
				log.Fatal("Failed to populate database: ", err)
			}
		}
	}

	log.Println("Database populated with initial data")
}

func (db *database) RawDB() *gorm.DB {
	return db.DB
}

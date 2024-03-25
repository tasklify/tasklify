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
	GetUsers() ([]User, error)
	GetFilteredUsers(userIDs []uint) ([]User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id uint) (*User, error)
	UpdateUser(user *User) error
	DeleteUserByID(id uint) error
	CreateUser(user *User) error
	GetSprintByProject(projectID uint) ([]Sprint, error)
	GetSprintByID(id uint) (*Sprint, error)
	CreateSprint(sprint *Sprint) error
	CreateUserStory(userStory *UserStory) error
	UpdateUserStory(userStory *UserStory) error
	AddUserStoryToSprint(sprintID uint, userStories []uint) (*Sprint, error)
	GetUserStoriesByProject(projectID uint) ([]UserStory, error)
	GetUserStoriesBySprint(sprintID uint) ([]UserStory, error)
	GetUserStoryByID(id uint) (*UserStory, error)
	UserStoryInThisProjectAlreadyExists(title string, projectID uint) bool
	GetTasksByUserStory(userStoryID uint) ([]Task, error)
	CreateTask(task *Task) error
	GetTaskByID(id uint) (*Task, error)
	GetProjectByID(id uint) (*Project, error)
	GetProjectRoles(userID uint, projectID uint) ([]ProjectRole, error)
	GetProjectHasUserByProjectAndUser(userID uint, projectID uint) (*ProjectHasUser, error)
	CreateProject(project *Project) (uint, error)
	AddDeveloperToProject(projectID uint, userID uint) error
	RemoveUsersNotInList(projectID uint, userIDs []uint) error
	GetUsersWithRoleOnProject(projectID uint, projectRole ProjectRole) ([]User, error)
	GetUsersNotOnProject(projectID uint) ([]User, error)
	ProjectWithTitleExists(title string, excludedProjectID *uint) bool
	RemoveUserFromProject(projectID uint, userID uint) error
	GetUserProjects(userID uint) ([]Project, error)
	UpdateProject(projectID uint, projectData Project) error
	CreateAcceptanceTest(acceptanceTest *AcceptanceTest) error
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
	err := db.AutoMigrate(&User{}, &Project{}, &UserStory{}, &Task{}, &AcceptanceTest{}, &ProjectHasUser{}, &Sprint{}, &WorkflowStep{})
	if err != nil {
		log.Fatal("Schema migration error: ", err)
	}

	log.Println("Database tables registered")
}

func (db *database) RawDB() *gorm.DB {
	return db.DB
}

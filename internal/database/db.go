package database

import (
	"fmt"
	"log"
	"sync"
	"tasklify/internal/config"
	"time"

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
	GetUserStoriesLoad(userStoryIDs []uint) (uint, error)
	GetUserStoriesByProject(projectID uint) ([]UserStory, error)
	GetUserStoriesBySprint(sprintID uint) ([]UserStory, error)
	GetUserStoryByID(id uint) (*UserStory, error)
	UserStoryInThisProjectAlreadyExists(title string, projectID uint) bool
	GetTasksByUserStory(userStoryID uint) ([]Task, error)
	CreateTask(task *Task) error
	UpdateTask(task *Task) error
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
	GetProjectWallPosts(projectID uint) ([]ProjectWallPost, error)
	AddProjectWallPost(post ProjectWallPost) error
	EditProjectWallPost(projectID uint, postID uint, body string) error
	DeleteProjectWallPost(projectID uint, postID uint) error
	GetProjectWallPostByID(postID uint) (*ProjectWallPost, error)
	CreateAcceptanceTest(acceptanceTest *AcceptanceTest) error
	StartWorkSession(userID, taskID uint) error
	ResumeWorkSession(sessionID uint) error
	EndWorkSession(sessionID uint) error
	GetTotalTimeSpentOnTask(taskID uint) (time.Duration, error)
	GetWorkSessionsForTask(taskID uint) ([]WorkSession, error)
	CloseOpenSessionsAtMidnight() error
	GetWorkSessionByID(sessionID uint) (*WorkSession, error)
	ChangeDuration(sessionID uint, duration time.Duration) error
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
	err := db.AutoMigrate(&User{}, &Project{}, &ProjectWallPost{}, &UserStory{}, &Task{}, &AcceptanceTest{}, &ProjectHasUser{}, &Sprint{}, &WorkflowStep{}, &WorkSession{})
	if err != nil {
		log.Fatal("Schema migration error: ", err)
	}

	log.Println("Database tables registered")
}

func (db *database) RawDB() *gorm.DB {
	return db.DB
}

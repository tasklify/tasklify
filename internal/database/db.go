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
	UpdateSprint(sprint *Sprint) error
	DeleteSprint(projectID uint, sprintID uint) error
	CreateUserStory(userStory *UserStory) error
	UpdateUserStory(userStory *UserStory) error
	AddUserStoryToSprint(sprintID uint, userStories []uint) (*Sprint, error)
	DeleteUserStory(userStoryID uint) error
	GetUserStoriesLoad(userStoryIDs []uint) (uint, error)
	GetUserStoriesByProject(projectID uint) ([]UserStory, error)
	GetUserStoriesBySprint(sprintID uint) ([]UserStory, error)
	GetUserStoryByID(id uint) (*UserStory, error)
	GetAcceptanceTestsByUserStory(userStoryID uint) ([]AcceptanceTest, error)
	DeleteAcceptanceTest(acceptanceTest *AcceptanceTest) error
	UserStoryInThisProjectAlreadyExists(title string, projectID uint) bool
	UserStoryInThisProjectAlreadyExistsEdit(title string, projectID uint, userStroyID uint) bool
	GetTasksByUserStory(userStoryID uint) ([]Task, error)
	CreateTask(task *Task) error
	GetTaskByID(id uint) (*Task, error)
	UpdateTask(task *Task) error
	DeleteTask(taskID uint) error
	GetProjectByID(id uint) (*Project, error)
	GetProjectRoles(userID uint, projectID uint) ([]ProjectRole, error)
	GetProjectHasUserByProjectAndUser(userID uint, projectID uint) (*ProjectHasUser, error)
	CreateProject(project *Project) (uint, error)
	DeleteProject(projectID uint) error
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
	GetAcceptanceTestByID(id uint) (*AcceptanceTest, error)
	UpdateAcceptanceTest(acceptanceTest *AcceptanceTest) error
	CreateWorkSession(session *WorkSession) error
	CreateWorkSessionToday(userID, taskID uint) error
	GetWorkSessionsForTask(taskID uint) ([]WorkSession, error)
	GetWorkSessionByID(sessionID uint) (*WorkSession, error)
	UpdateWorkSession(session *WorkSession) error
	DeleteWorkSession(sessionID uint) error
	GetUserStoryComments(userStoryID uint) ([]UserStoryComment, error)
	GetUserStoryCommentByID(commentID uint) (*UserStoryComment, error)
	AddUserStoryComment(comment UserStoryComment) error
	EditUserStoryComment(userStoryID uint, commentID uint, body string) error
	DeleteUserStoryComment(userStoryID uint, commentID uint) error
	GetUserTasks(userID uint) ([]Task, error)
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
	err := db.AutoMigrate(&User{}, &Project{}, &ProjectWallPost{}, &UserStory{}, &UserStoryComment{}, &Task{}, &AcceptanceTest{}, &ProjectHasUser{}, &Sprint{}, &WorkflowStep{}, &WorkSession{})
	if err != nil {
		log.Fatal("Schema migration error: ", err)
	}

	log.Println("Database tables registered")
}

func (db *database) RawDB() *gorm.DB {
	return db.DB
}

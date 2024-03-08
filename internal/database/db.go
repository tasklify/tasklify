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
	GetSystemRole(systemRoleName string) (*SystemRole, error)
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
		registerEnums(databaseClient)
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

func registerEnums(db *database) {
	// Create ENUMs
	if err := db.Exec(`
        DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum') THEN
                CREATE TYPE status_enum AS ENUM ('StatusTodo', 'StatusInProgress', 'StatusDone');
            END IF;

            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'priority_enum') THEN
                CREATE TYPE priority_enum AS ENUM ('PriorityMustHave', 'PriorityCouldHave', 'PriorityShouldHave', 'PriorityWontHaveThisTime');
            END IF;
        END$$;
    `).Error; err != nil {
		log.Fatal(err)
	}
}

func registerTables(db *database) {
	// Migrate the schema
	err := db.AutoMigrate(&User{}, &SystemRole{}, &ProjectRole{}, &Project{}, &UserStory{}, &Task{}, &ProjectHasUser{}, &Sprint{}, &WorkflowStep{})
	if err != nil {
		log.Fatal("Schema migration error: ", err)
	}

	log.Println("Database tables registered")
}

func populateDatabase(db *database) {
	var count int64

	// Create System Roles
	if db.Model(&SystemRole{}).Count(&count); count == 0 {
		for _, role := range systemRoles {
			if err := db.CreateSystemRole(&role); err != nil {
				log.Fatal("Failed to populate database: ", err)
			}
		}
	}

	// Create Project Roles
	if db.Model(&ProjectRole{}).Count(&count); count == 0 {
		for _, role := range projectRoles {
			if err := db.CreateProjectRole(&role); err != nil {
				log.Fatal("Failed to populate database: ", err)
			}
		}
	}

	// Create Users
	if db.Model(&User{}).Count(&count); count == 0 {
		for _, user := range users {
			// Associate user with system role
			var systemRole SystemRole
			if err := db.Where(&user.SystemRole).First(&systemRole).Error; err != nil {
				log.Fatal("Failed to retrieve SystemRole: ", err)
			}

			user.SystemRoleID = systemRole.ID
			user.SystemRole = SystemRole{} // Unset SystemRole, because we only need SystemRoleID to create user

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

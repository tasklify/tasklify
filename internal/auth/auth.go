package auth

import (
	"log"
	"tasklify/internal/config"
	"tasklify/internal/database"

	"github.com/aws/smithy-go/ptr"
)

const (
	RootUserID uint = 1
)

func InitialUsers() {
	err := CreateUser(nil, ptr.Uint(RootUserID), config.GetConfig().Admin.Username, config.GetConfig().Admin.Password, config.GetConfig().Admin.Password, "admin", "admin", "admin@mail.com", database.SystemRoleAdmin.Val)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Created initial users")
}

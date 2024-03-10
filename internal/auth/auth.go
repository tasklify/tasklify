package auth

import (
	"log"
	"tasklify/internal/config"
	"tasklify/internal/database"

	"github.com/aws/smithy-go/ptr"
)

func InitialUsers() {
	err := CreateUser(ptr.Uint(9999), config.GetConfig().Admin.Username, config.GetConfig().Admin.Password, "admin", "admin", "admin@mail.com", database.SystemRoleAdmin.Val)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Created initial users")
}

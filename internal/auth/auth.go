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
	err := CreateUser(nil, ptr.Uint(RootUserID), config.GetConfig().Admin.Username, config.GetConfig().Admin.Password, "admin", "admin", "admin@mail.com", database.SystemRoleAdmin.Val)
	if err != nil {
		log.Panic(err)
	}

	err = CreateUser(nil, ptr.Uint(2), "test1", "test1@mail.com", "Testni1", "Uporabnik1", "test1@mail.com", database.SystemRoleUser.Val)
	if err != nil {
		log.Panic(err)
	}

	err = CreateUser(nil, ptr.Uint(3), "test2", "test1@mail.com", "Testni2", "Uporabnik2", "test2@mail.com", database.SystemRoleUser.Val)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Created initial users")
}

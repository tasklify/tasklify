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

	err = CreateUser(ptr.Uint(8888), "test1", "test1", "Testni1", "Uporabnik1", "test1@mail.com", database.SystemRoleUser.Val)
	if err != nil {
		log.Panic(err)
	}

	err = CreateUser(ptr.Uint(7777), "test2", "test2", "Testni2", "Uporabnik2", "test2@mail.com", database.SystemRoleUser.Val)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Created initial users")
}

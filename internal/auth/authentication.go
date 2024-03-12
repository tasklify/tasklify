package auth

import (
	"errors"
	"fmt"
	"log"
	"tasklify/internal/database"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gookit/goutil/dump"
)

func AuthenticateUser(username, password string) (uint, error) {
	loginTime := time.Now()

	user, err := database.GetDatabase().GetUserByUsername(username)
	if err != nil {
		return 0, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.Password)
	if err != nil {
		log.Println(err)
		return 0, errors.New("error when checking passwords")
	}

	if !match {
		return 0, errors.New("no matching username and password")
	}

	var userLastLogin = &database.User{}
	userLastLogin.ID = user.ID
	userLastLogin.LastLogin = &loginTime
	err = database.GetDatabase().UpdateUser(userLastLogin)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func CreateUser(ID *uint, username, password, firstName, lastName, email, systemRoleName string) error {
	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	systemRole := database.SystemRoles.Parse(systemRoleName)
	if systemRole == (&database.SystemRole{}) {
		return errors.New("system role not found")
	}

	var user = &database.User{
		Username:   username,
		Password:   passwordHash,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		SystemRole: *systemRole,
	}

	if ID != nil {
		user.ID = *ID
	}

	dump.P(user)

	return database.GetDatabase().UpdateUser(user)
}

func UpdateUser(issuerUserID, issuerPassword string, userID uint, username, password, firstName, lastName, email, systemRoleName *string) error {
	ok, err := AuthenticateUser(issuerUserID, issuerPassword)
	if err != nil {
		return err
	}

	issuerUser, err := database.GetDatabase().GetUserByUsername(issuerUserID)
	if err != nil {
		return err
	}

	if ok == 0 {
		return errors.New("you are not authenticated")
	}

	var user = &database.User{}
	user.ID = userID

	if username != nil {
		user.Username = *username
	}

	if password != nil {
		passwordHash, err := argon2id.CreateHash(*password, argon2id.DefaultParams)
		if err != nil {
			return err
		}

		user.Username = passwordHash
	}

	if firstName != nil {
		user.FirstName = *firstName
	}

	if lastName != nil {
		user.LastName = *lastName
	}

	if email != nil {
		user.Email = *email
	}

	if systemRoleName != nil {
		err = GetAuthorization().HasSystemPermission(issuerUser.SystemRole, "/system/user/system-role", ActionUpdate)
		if err != nil {
			return err
		}

		systemRole := database.SystemRoles.Parse(*systemRoleName)
		if systemRole == nil {
			return fmt.Errorf("'%s' is not a valid system role", *systemRoleName)
		}

		user.SystemRole = *systemRole
	}

	return database.GetDatabase().UpdateUser(user)
}

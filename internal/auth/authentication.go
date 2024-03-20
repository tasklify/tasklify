package auth

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"log"
	"tasklify/internal/config"
	"tasklify/internal/database"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gookit/goutil/dump"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func AuthenticateUser(username, password string) (uint, error) {
	loginTime := time.Now()

	err := checkPasswordRequirements(password)
	if err != nil {
		return 0, err
	}

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
	userLastLogin = user
	userLastLogin.LastLogin = &loginTime
	err = database.GetDatabase().UpdateUser(userLastLogin)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func CreateUser(ID *uint, username, password, firstName, lastName, email, systemRoleName string) error {
	err := checkPasswordRequirements(password)
	if err != nil {
		return err
	}

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

	if config.GetConfig().Debug {
		dump.P(user)
	}

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
		err := checkPasswordRequirements(*password)
		if err != nil {
			return err
		}

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

// CheckPasswordRequirements verifies the password against the specified rules.
func checkPasswordRequirements(password string) error {
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}
	if len(password) > 128 {
		return errors.New("password must not be longer than 128 characters")
	}

	err := GetCommonPasswordList().IsCommon(password)
	if err != nil {
		return err
	}

	return nil
}

func GenerateTOTP(userID uint) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "tasklify",
		AccountName: fmt.Sprint(userID),
	})
}

func DisplayTOTP(key *otp.Key) (string, *bytes.Buffer, error) {
	var buf *bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", nil, err
	}

	err = png.Encode(buf, img)
	if err != nil {
		return "", nil, err
	}

	return key.Secret(), buf, nil
}

func VerifyTOTP(passcode, secret string) error {
	valid := totp.Validate(passcode, secret)

	if !valid {
		return fmt.Errorf("passcode not valid")
	}

	return nil
}

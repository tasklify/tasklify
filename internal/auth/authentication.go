package auth

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"log"
	"tasklify/internal/database"
	"time"

	"github.com/alexedwards/argon2id"
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

func CreateUser(issuerUserID *uint, userID *uint, username, password, firstName, lastName, email, systemRoleName string) error {
	if issuerUserID != nil {
		issuerUser, err := database.GetDatabase().GetUserByID(*issuerUserID)
		if err != nil {
			return err
		}

		if issuerUser.SystemRole != database.SystemRoleAdmin {
			return fmt.Errorf("system_role admin required for this action")
		}
	}

	err := checkPasswordRequirements(password)
	if err != nil {
		return err
	}

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	systemRole := database.SystemRoles.Parse(systemRoleName)
	if systemRole == nil {
		return errors.New("invalid system role")
	}

	var user = &database.User{
		Username:   username,
		Password:   passwordHash,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		SystemRole: *systemRole,
	}

	if userID != nil {
		user.ID = *userID
	}

	// if config.GetConfig().Debug {
	// 	dump.P(user)
	// }

	return database.GetDatabase().UpdateUser(user)
}

func DeleteUser(issuerUserID uint, issuerPassword string, userID uint) error {
	issuerUser, err := database.GetDatabase().GetUserByID(issuerUserID)
	if err != nil {
		return err
	}

	_, err = AuthenticateUser(issuerUser.Username, issuerPassword)
	if err != nil {
		return err
	}

	if issuerUser.SystemRole != database.SystemRoleAdmin {
		return fmt.Errorf("system_role admin required for this action")
	}

	// if config.GetConfig().Debug {
	// 	dump.P(user)
	// }

	return database.GetDatabase().DeleteUserByID(userID)
}

func UpdateUser(issuerUserID uint, issuerPassword string, userID uint, username, password, firstName, lastName, email, systemRoleName *string) error {
	issuerUser, err := database.GetDatabase().GetUserByID(issuerUserID)
	if err != nil {
		return err
	}

	ok, err := AuthenticateUser(issuerUser.Username, issuerPassword)
	if err != nil {
		return err
	}

	if ok == 0 {
		return errors.New("you are not authenticated")
	}

	user, err := database.GetDatabase().GetUserByID(userID)
	if err != nil {
		return err
	}

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

		user.Password = passwordHash
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
	passwordLen := len(password)

	if passwordLen < 12 {
		return fmt.Errorf("password must be at least 12 characters long, currently %d", passwordLen)
	}
	if passwordLen > 128 {
		return fmt.Errorf("password must not be longer than 128 characters, currently %d", passwordLen)
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

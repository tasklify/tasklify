package auth

import (
	"fmt"
	"log"
	"sync"

	dumbpassword "github.com/theifedayo/go-dumb-password"
)

type CommonPasswordList interface {
	IsCommon(pasword string) error
}

type commonPasswordList struct {
	*dumbpassword.DumbPasswordValidator
}

var (
	onceCommonPasswordList sync.Once

	commonPasswordListClient *commonPasswordList
)

func GetCommonPasswordList() CommonPasswordList {

	onceCommonPasswordList.Do(func() { // <-- atomic, does not allow repeating
		commonPasswordListClient = loadCommonPasswordList()
	})

	return commonPasswordListClient
}

func loadCommonPasswordList() *commonPasswordList {
	passwordListPath := "configs/common_password_list.txt"
	validator, err := dumbpassword.DPValidator(passwordListPath)
	if err != nil {
		panic(err)
	}

	log.Println("Common Password List loaded")

	return &commonPasswordList{DumbPasswordValidator: validator}
}

func (c *commonPasswordList) IsCommon(pasword string) error {
	notCommon := c.Validate(pasword)

	if !notCommon {
		return fmt.Errorf("password is in the most common password list of 1 million passwords")
	}

	return nil
}

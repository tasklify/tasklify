package auth

import (
	"errors"
	"log"
	"sync"

	"github.com/casbin/casbin/v2"
)

type Authorization interface {
	HasPermission(projectGroup, object, action string) error
}

type authorization struct {
	*casbin.Enforcer
}

var (
	onceAuthorization sync.Once

	authorizationClient *authorization
)

func GetAuthorization() Authorization {

	onceAuthorization.Do(func() { // <-- atomic, does not allow repeating
		authorizationClient = connectAuthorization()
	})

	return authorizationClient
}

func connectAuthorization() *authorization {
	// a, _ := gormadapter.NewAdapterByDB(database.GetDatabase().RawDB())
	enforcer, err := casbin.NewEnforcer("./rbac/model.conf", "./rbac/policy.csv" /*, a */)
	if err != nil {
		log.Panic(err)
	}

	// Load the policy
	enforcer.LoadPolicy()

	log.Println("Authorization connected")

	return &authorization{Enforcer: enforcer}
}

func (a *authorization) HasPermission(groupName, object, action string) error {
	ok, err := a.Enforce(groupName, object, action)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("not authorized")
	}

	return nil
}

package auth

import (
	"errors"
	"log"
	"sync"
	"tasklify/internal/database"

	"github.com/casbin/casbin/v2"
)

type Authorization interface {
	HasSystemPermission(systemRole database.SystemRole, object string, action Action) error
	HasProjectPermission(systemRole database.ProjectRole, object string, action Action) error
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
	enforcer, err := casbin.NewEnforcer("./configs/rbac_model.conf", "./configs/rbac_policy.csv" /*, a */)
	if err != nil {
		log.Panic(err)
	}

	// Load the policy
	enforcer.LoadPolicy()

	log.Println("Authorization connected")

	return &authorization{Enforcer: enforcer}
}

func (a *authorization) HasSystemPermission(systemRole database.SystemRole, object string, action Action) error {
	ok, err := a.Enforce(systemRole.Val, object, action.Val)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("not authorized")
	}

	return nil
}

func (a *authorization) HasProjectPermission(systemRole database.ProjectRole, object string, action Action) error {
	ok, err := a.Enforce(systemRole.Val, object, action.Val)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("not authorized")
	}

	return nil
}

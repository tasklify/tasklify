package auth

import (
	"github.com/casbin/casbin/v2"
)

func Casbin() {
	// a, _ := gormadapter.NewAdapterByDB(database.GetDatabase().RawDB())
	e, _ := casbin.NewEnforcer("./rbac/model.conf", "./rbac/policy.csv" /*, a */)

	// Load the policy
	e.LoadPolicy()

	e.Enforce()
}

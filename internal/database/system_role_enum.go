package database

import (
	"database/sql/driver"
	"fmt"

	"github.com/orsinium-labs/enum"
)

type SystemRole enum.Member[string]

var (
	e               = enum.NewBuilder[string, SystemRole]()
	SystemRoleAdmin = e.Add(SystemRole{"system_admin"})
	SystemRoleUser  = e.Add(SystemRole{"system_user"})
	SystemRoles     = e.Enum()
)

func (role SystemRole) Value() (driver.Value, error) {
	fmt.Println("here1")
	return driver.Value(SystemRoles.Value(role)), nil
}

func (role *SystemRole) Scan(value interface{}) error {
	fmt.Println("here2")
	*role = *SystemRoles.Parse(value.(string))
	return nil
}

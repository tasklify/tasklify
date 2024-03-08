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

// CustomType embeds SystemRole and implements driver.Valuer
type SystemRoleGORM struct {
	*SystemRole
}

func (role *SystemRoleGORM) Value() (driver.Value, error) {
	fmt.Println("here1")
	return driver.Value(SystemRoles.Value(*role.SystemRole)), nil
}

func (role *SystemRoleGORM) Scan(value interface{}) error {
	fmt.Println("here2")
	systemRole := SystemRoles.Parse(value.(string))

	*role = SystemRoleGORM{SystemRole: systemRole}
	return nil
}

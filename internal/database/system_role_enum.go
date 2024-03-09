package database

import (
	"database/sql/driver"
	"tasklify/pkg/enum"
)

type SystemRole enum.Member[string]

var (
	SystemRoleEnum  = enum.NewBuilder[string, SystemRole]()
	SystemRoleAdmin = SystemRoleEnum.Add(SystemRole{"system_admin"})
	SystemRoleUser  = SystemRoleEnum.Add(SystemRole{"system_user"})
	SystemRoles     = SystemRoleEnum.Enum()
)

func (role SystemRole) Value() (driver.Value, error) {
	if role == (SystemRole{}) {
		return nil, nil
	}
	return SystemRoles.WrappedValue(role), nil
}

func (role *SystemRole) Scan(value interface{}) error {
	parsedRole := SystemRoles.Parse(value.(string))
	*role = *parsedRole

	return nil
}

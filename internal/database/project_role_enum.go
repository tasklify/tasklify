package database

import (
	"database/sql/driver"
	"tasklify/pkg/enum"
)

type ProjectRole enum.Member[string]

var (
	ProjectRoleEnum      = enum.NewBuilder[string, ProjectRole]()
	ProjectRoleManager   = ProjectRoleEnum.Add(ProjectRole{"project_manager"})
	ProjectRoleMaster    = ProjectRoleEnum.Add(ProjectRole{"project_master"})
	ProjectRoleDeveloper = ProjectRoleEnum.Add(ProjectRole{"project_developer"})
	ProjectRoles         = ProjectRoleEnum.Enum()
)

func (role ProjectRole) Value() (driver.Value, error) {
	if role == (ProjectRole{}) {
		return nil, nil
	}
	return ProjectRoles.WrappedValue(role), nil
}

func (role *ProjectRole) Scan(value interface{}) error {
	parsedRole := ProjectRoles.Parse(value.(string))
	*role = *parsedRole

	return nil
}

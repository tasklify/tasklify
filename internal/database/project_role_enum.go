package database

import (
	"database/sql/driver"
	"tasklify/third_party/enum"
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

func (role ProjectRole) GetLabel() string {
	switch role {
	case ProjectRoleManager:
		return "Project manager"
	case ProjectRoleMaster:
		return "Project master"
	case ProjectRoleDeveloper:
		return "Project developer"
	default:
		return "Unknown"
	}
}

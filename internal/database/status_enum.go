package database

import (
	"database/sql/driver"
	"tasklify/pkg/enum"
)

type Status enum.Member[string]

var (
	StatusEnum       = enum.NewBuilder[string, Status]()
	StatusTodo       = StatusEnum.Add(Status{"status_todo"})
	StatusInProgress = StatusEnum.Add(Status{"status_in_progress"})
	StatusDone       = StatusEnum.Add(Status{"status_done"})
	Statuses         = StatusEnum.Enum()
)

func (status Status) Value() (driver.Value, error) {
	if status == (Status{}) {
		return nil, nil
	}
	return Statuses.WrappedValue(status), nil
}

func (status *Status) Scan(value interface{}) error {
	parsedStatus := Statuses.Parse(value.(string))
	*status = *parsedStatus

	return nil
}

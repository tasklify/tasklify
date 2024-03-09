package database

import (
	"database/sql/driver"
	"tasklify/pkg/enum"
)

type Priority enum.Member[string]

var (
	PriorityEnum             = enum.NewBuilder[string, Priority]()
	PriorityMustHave         = PriorityEnum.Add(Priority{"priority_must_have"})
	PriorityCouldHave        = PriorityEnum.Add(Priority{"priority_could_have"})
	PriorityShouldHave       = PriorityEnum.Add(Priority{"priority_should_have"})
	PriorityWontHaveThisTime = PriorityEnum.Add(Priority{"priority_wont_have_this_time"})
	Priorities               = PriorityEnum.Enum()
)

func (priority Priority) Value() (driver.Value, error) {
	if priority == (Priority{}) {
		return nil, nil
	}
	return Priorities.WrappedValue(priority), nil
}

func (priority *Priority) Scan(value interface{}) error {
	parsedPriority := Priorities.Parse(value.(string))
	*priority = *parsedPriority

	return nil
}

package config

import "fmt"

type EnvironmentEnum struct {
	slug string
}

func (o EnvironmentEnum) String() string {
	return o.slug
}

var (
	Invalid = EnvironmentEnum{"invalid"}
	Prod    = EnvironmentEnum{"prod"}
	Dev     = EnvironmentEnum{"dev"}
)

func EnvironmentEnums() []EnvironmentEnum {
	return []EnvironmentEnum{Prod, Dev}
}

func EnvironmentEnumFromString(s string) (EnvironmentEnum, error) {
	switch s {
	case Prod.slug:
		return Prod, nil
	case Dev.slug:
		return Dev, nil
	}

	return Invalid, fmt.Errorf("no valid environmanrt found in input: %s", s)
}

// Only for parsers
func EnvironmentEnumFromStringParser(v string) (interface{}, error) {
	env, err := EnvironmentEnumFromString(v)
	if err != nil {
		return nil, err
	}
	return env, nil
}

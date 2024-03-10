package config

import (
	"fmt"
	"tasklify/third_party/enum"
)

type Environment enum.Member[string]

var (
	e               = enum.NewBuilder[string, Environment]()
	EnvironmentProd = e.Add(Environment{"prod"})
	EnvironmentDev  = e.Add(Environment{"dev"})
	Environments    = e.Enum()
)

// Only for parsers
func environmentParser(v string) (interface{}, error) {
	env := Environments.Parse(v)
	if env == nil {
		return nil, fmt.Errorf("'%s' is not a valid environment enum", v)
	}
	return *env, nil
}

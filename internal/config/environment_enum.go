package config

import (
	"fmt"

	"github.com/orsinium-labs/enum"
)

type Environment enum.Member[string]

var (
	e            = enum.NewBuilder[string, Environment]()
	Prod         = e.Add(Environment{"prod"})
	Dev          = e.Add(Environment{"dev"})
	Environments = e.Enum()
)

// Only for parsers
func environmentParser(v string) (interface{}, error) {
	env := Environments.Parse(v)
	if env == nil {
		return nil, fmt.Errorf("'%s' is not a valid enviranment enum", v)
	}
	return *env, nil
}

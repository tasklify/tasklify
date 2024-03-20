package config

import (
	"log"
	"reflect"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/gookit/goutil/dump"
)

type Config struct {
	Debug       bool        `env:"TASKLIFY_DEBUG" envDefault:"false"`
	Environment Environment `env:"TASKLIFY_ENVIRONMENT" envDefault:"prod"`
	Host        string      `env:"TASKLIFY_HOST" envDefault:"0.0.0.0"`
	Port        string      `env:"TASKLIFY_PORT" envDefault:"8080"`
	Database    Database
	Auth        Auth
	Admin       Admin
}

type Database struct {
	Host     string `env:"TASKLIFY_DATABASE_HOST"`
	Port     string `env:"TASKLIFY_DATABASE_PORT" envDefault:"5432"`
	DbName   string `env:"TASKLIFY_DATABASE_NAME"`
	User     string `env:"TASKLIFY_DATABASE_USER"`
	Password string `env:"TASKLIFY_DATABASE_PASSWORD"`
}

type Auth struct {
	SessionHashKey  string `env:"TASKLIFY_AUTH_SESSION_HASH_KEY"`
	SessionBlockKey string `env:"TASKLIFY_AUTH_SESSION_ENCRYPTION_KEY"`
	SymcryptKey     string `env:"TASKLIFY_AUTH_SYMCRYPT_KEY"`
}

type Admin struct {
	Username string `env:"TASKLIFY_ADMIN_USERNAME"`
	Password string `env:"TASKLIFY_ADMIN_PASSWORD"`
}

var (
	onceConfig sync.Once

	config *Config
)

func GetConfig() *Config {

	onceConfig.Do(func() { // <-- atomic, does not allow repeating
		config = loadConfig()
	})

	return config
}

func loadConfig() *Config {
	config := &Config{}
	envOptions := env.Options{RequiredIfNoDef: true,
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(Environment{}): environmentParser,
		},
	}

	// Load env vars
	err := env.ParseWithOptions(config, envOptions)
	if err != nil {
		log.Panic(err)
	}

	if config.Debug {
		log.Println("Running in DEBUG mode")
		dump.P(config)
	}

	return config
}

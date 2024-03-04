package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/gookit/goutil/dump"
)

type Environment string

const (
	PROD Environment = "PROD"
	DEV  Environment = "DEV"
)

type Config struct {
	Debug    bool   `env:"TASKLIFY_DEBUG" envDefault:"false"`
	Host     string `env:"TASKLIFY_HOST" envDefault:"0.0.0.0"`
	Port     string `env:"TASKLIFY_PORT" envDefault:"8080"`
	Database Database
	Auth     Auth
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
	envOptions := env.Options{RequiredIfNoDef: true}

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

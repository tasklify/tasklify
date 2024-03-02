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
	Debug    bool   `env:"DEBUG" envDefault:"false"`
	Host     string `env:"HOST" envDefault:"0.0.0.0"`
	Port     string `env:"PORT" envDefault:"8080"`
	Database Database
}

type Database struct {
	Host     string `env:"TASKLIFY_DATABASE_HOST"`
	Port     string `env:"TASKLIFY_DATABASE_HOST" envDefault:"5432"`
	DbName   string `env:"TASKLIFY_DATABASE_NAME"`
	User     string `env:"TASKLIFY_DATABASE_USER"`
	Password string `env:"TASKLIFY_DATABASE_PASSWORD"`
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
		log.Fatal(err)
	}

	if config.Debug {
		log.Println("Running in DEBUG mode")
		dump.P(config)
	}

	return config
}

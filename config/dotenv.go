package config

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

var LoadEnvFile = func(envPath string) error {
	return godotenv.Load(envPath)
}

func Env(name string) string {
	var envPath string
	if testing.Testing() {
		envPath = "../.env"
	} else {
		envPath = ".env"
	}

	err := LoadEnvFile(envPath)
	if err != nil {
		if !testing.Testing() {
			log.Printf("Error loading .env file: %s", err)
		}
	}
	return os.Getenv(name)
}

func EnvAsInt(name string, defaultVal int) int {
	valueStr := Env(name)
	if valueStr == "" {
		return defaultVal
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Could not parse %s as integer, using default value %d", name, defaultVal)
		return defaultVal
	}

	return value
}

package config

import (
	"log"
	"os"
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
			log.Printf("Помилка завантаження файлу .env: %s", err)
		}
	}
	return os.Getenv(name)
}

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variable")
	}
}

func LoadEnvTest() {
	if err := godotenv.Load(".env.test"); err != nil {
		fmt.Println("No .env.test file found, using environment variable")
	}
}

func Get(key string) string {
	return os.Getenv(key)
}

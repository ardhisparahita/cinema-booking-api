package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variable")
	}
}

func LoadEnvTest() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("unable to get working directory:", err)
		return
	}
	for {
		envPath := filepath.Join(dir, ".env.test")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				fmt.Println("failed to load .env.test:", err)
			} else {
				fmt.Println("loaded:", envPath)
			}
			return
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			fmt.Println("no .env.test file found, using env variable")
			return
		}

		dir = parent
	}
}

func Get(key string) string {
	return os.Getenv(key)
}

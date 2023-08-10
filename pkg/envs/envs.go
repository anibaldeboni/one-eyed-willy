package envs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func IsDev() bool {
	return Current() == "development"
}

func IsProd() bool {
	return Current() == "production"
}

func Current() string {
	currentEnv := Get("ENVIRONMENT")
	if currentEnv == "" {
		currentEnv = "development"
	}
	return currentEnv
}

func Load() error {
	envFile := fmt.Sprintf(".env.%s", Current())

	_, err := os.Stat(envFile)

	if err == nil {
		err = godotenv.Load(os.ExpandEnv(envFile))
		if err != nil {
			return fmt.Errorf("error initializing app: %v", err)
		}
	}
	return nil
}

func Get(key string) string {
	return os.Getenv(key)
}

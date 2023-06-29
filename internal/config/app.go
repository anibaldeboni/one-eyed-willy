package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AppConfig struct {
	Environment string `json:"environment"`
	AppPort     string `json:"app_port"`
	BaseURL     string `json:"base_url"`

	Validator  echo.Validator        `json:"-"`
	CORSConfig middleware.CORSConfig `json:"-"`

	LogLevel string `json:"log_level"`
}

type AppValidator struct {
	validator *validator.Validate
}

func (cv *AppValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func InitAppConfig() (*AppConfig, error) {
	var err error

	currentEnv := os.Getenv("ENVIRONMENT")
	if currentEnv == "" {
		currentEnv = "development"
	}

	envFile := fmt.Sprintf(".env.%s", currentEnv)

	_, err = os.Stat(envFile)

	if err == nil {
		err = godotenv.Load(os.ExpandEnv(envFile))
		if err != nil {
			return nil, fmt.Errorf("error initializing app: %v", err)
		}
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &AppConfig{
		Environment: currentEnv,
		AppPort:     appPort,
		BaseURL:     os.Getenv("BASE_URL"),
		Validator:   &AppValidator{validator: validator.New()},
		CORSConfig:  middleware.DefaultCORSConfig,

		LogLevel: os.Getenv("LOG_LEVEL"),
	}, nil
}

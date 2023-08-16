package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/one-eyed-willy/pkg/envs"
)

type AppConfig struct {
	Environment string `json:"environment"`
	AppPort     string `json:"app_port"`
	BaseURL     string `json:"base_url"`

	Validator  echo.Validator        `json:"-"`
	CORSConfig middleware.CORSConfig `json:"-"`

	LogLevel    string `json:"log_level"`
	MaxFileSize string `json:"max_file_size"`
}

type AppValidator struct {
	validator *validator.Validate
}

func (cv *AppValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func InitAppConfig() *AppConfig {
	appPort := envs.Get("PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &AppConfig{
		Environment: envs.Current(),
		AppPort:     appPort,
		BaseURL:     os.Getenv("BASE_URL"),
		Validator:   &AppValidator{validator: validator.New()},
		CORSConfig:  middleware.DefaultCORSConfig,

		LogLevel:    envs.Get("LOG_LEVEL"),
		MaxFileSize: "5M",
	}
}

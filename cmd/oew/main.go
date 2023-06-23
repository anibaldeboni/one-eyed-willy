package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/one-eyed-willy/internal/config"
	"github.com/one-eyed-willy/pkg/logger"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

// @title One-Eyed-Willy REST pdf generation API
// @version 1.0
// @description This documentation for One-Eyed-Willy pdf generator.
// @termsOfService http://swagger.io/terms/

// @contact.name Aníbal Deboni Neto
// @contact.email anibaldeboni@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /pdf
// @schemes http https
// @produces application/json application/octet-stream
// @consumes application/json
func main() {
	conf, _ := config.InitAppConfig()

	logger.InitWithOptions(logger.WithConfigLevel(conf.LogLevel))
	if logger.Log() != nil {
		defer logger.Log().Sync()
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(config.GetEchoLogConfig(conf)))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.HideBanner = true
	e.Validator = conf.Validator

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello there! This is One-Eyed-Willy pdf generator",
		})
	})

	e.GET("/docs/*", echoSwagger.WrapHandler)

	// Start server
	go func() {
		if err := e.Start(":" + conf.AppPort); err != nil && err != http.ErrServerClosed {
			logger.Log().Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 5 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/one-eyed-willy/docs"
	"github.com/one-eyed-willy/internal/config"
	"github.com/one-eyed-willy/internal/handler"
	"github.com/one-eyed-willy/internal/router"
	"github.com/one-eyed-willy/pkg/logger"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			One-Eyed-Willy REST pdf generation API
//	@version		1.0
//	@description	This documentation for One-Eyed-Willy pdf generator.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Aníbal Deboni Neto
//	@contact.email	anibaldeboni@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/pdf
// @schemes	http https
// @produces	application/json application/octet-stream
// @consumes	application/json
func main() {
	conf, _ := config.InitAppConfig()
	r := router.New(conf)
	h := handler.New()

	r.GET("/", h.IndexView)

	r.GET("/docs/*", echoSwagger.WrapHandler)

	h.Register(r.Group(conf.BaseURL))

	go func() {
		if err := r.Start(":" + conf.AppPort); err != nil && err != http.ErrServerClosed {
			logger.Log().Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 5 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		logger.Log().Info("shutting down server, good bye!")
		cancel()
	}()
	if err := r.Shutdown(ctx); err != nil {
		r.Logger.Fatal(err)
	}
}

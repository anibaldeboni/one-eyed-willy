package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/one-eyed-willy/docs"
	"github.com/one-eyed-willy/internal/handler"
	"github.com/one-eyed-willy/internal/web"
	"github.com/one-eyed-willy/pkg/envs"
	"github.com/one-eyed-willy/pkg/logger"
)

//	@title  One-Eyed-Willy pdf generation API
//	@version  1.0
//	@description  This documentation for One-Eyed-Willy pdf generator.
//	@termsOfService  http://swagger.io/terms/

//	@contact.name  Aníbal Deboni Neto
//	@contact.email  anibaldeboni@gmail.com

//	@license.name  Apache 2.0
//	@license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /pdf
// @schemes  http https
// @produces  application/json application/octet-stream
// @consumes  application/json
func main() {
	envs.Load()
	logger := logger.New()
	defer func() {
		_ = logger.Sync()
	}()
	w := web.New(logger)
	h := handler.New(w, logger)
	defer h.PdfRender.Cancel()

	go func() {
		if err := w.Start(":" + appPort()); err != nil && err != http.ErrServerClosed {
			h.PdfRender.Cancel()
			logger.Fatal(fmt.Sprintf("shutting down the server: %v", err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 5 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		fmt.Println("\nshutting down the server, good bye!")
		cancel()
	}()
	if err := w.Shutdown(ctx); err != nil {
		w.Logger.Fatal(err)
	}
}

func appPort() string {
	var port = envs.Get("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/one-eyed-willy/docs"
	"github.com/one-eyed-willy/internal/config"
	"github.com/one-eyed-willy/internal/handler"
	"github.com/one-eyed-willy/internal/web"
	"github.com/one-eyed-willy/pkg/logger"
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
	w := web.New(conf)
	h, err := handler.New()
	if err != nil {
		logger.Log().Fatal(err)
	}
	defer h.CancelPdfRenderContext()

	h.Register(w, conf)

	go func() {
		if err := w.Start(":" + conf.AppPort); err != nil && err != http.ErrServerClosed {
			logger.Log().Fatalf("shutting down the server: %v", err)
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

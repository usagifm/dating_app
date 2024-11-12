package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/middleware/request"
	v1 "github.com/usagifm/dating-app/src/v1"
)

//	@title			Dating App
//	@version		1.0
//	@description	Dating App for Dealls

//	@contact.name	@usagifm
//	@contact.url	https://example.id

//	@host		dating-app.taktix.co.id
//	@BasePath	/

func main() {
	initCtx := context.Background()
	if err := app.Init(initCtx); err != nil {
		panic(err)
	}

	startService(initCtx)
}

func startService(ctx context.Context) {
	address := fmt.Sprintf(":%d", app.Config().BindAddress)
	logger.GetLogger(ctx).Infof("Starting payment api service on %s", address)

	r := chi.NewRouter()
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(request.RequestIDContext(request.DefaultGenerator))
	r.Use(request.RequestAttributesContext)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	logger.GetLogger(context.Background()).Printf("1")
	deps := v1.Dependencies(ctx)
	logger.GetLogger(context.Background()).Printf("2")
	v1.Router(r, deps)
	logger.GetLogger(context.Background()).Printf("3")

	server := &http.Server{Addr: address, Handler: r}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.GetLogger(ctx).Info("Server Started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger(ctx).Fatalf("Server err: %s\n", err)
		}
	}()

	<-done
	logger.GetLogger(ctx).Info("Stopping Server")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		logger.GetLogger(ctx).Fatalf("Server Shutdown Failed:%+v", err)
	}

	logger.GetLogger(ctx).Info("Server Exited Properly")
}

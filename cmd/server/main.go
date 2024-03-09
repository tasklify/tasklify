package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tasklify/internal/auth"
	"tasklify/internal/config"
	"tasklify/internal/database"
	"tasklify/internal/router"
	"time"

	_ "go.uber.org/automaxprocs" // Automatically set GOMAXPROCS to match Linux container CPU quota
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	config := config.GetConfig()
	database.GetDatabase(config)
	auth.GetSession(config)
	auth.GetAuthorization()

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: router.Router(),
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			log.Panicf("error starting server: %s\n", err)
		}
	}()

	logger.Info("Server started", slog.String("port", config.Port))
	<-killSig

	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}

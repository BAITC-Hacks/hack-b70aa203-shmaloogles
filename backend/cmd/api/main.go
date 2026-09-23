package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shmaloogles/business-task-platform/backend/internal/database"
	"github.com/shmaloogles/business-task-platform/backend/internal/server"
	"github.com/shmaloogles/business-task-platform/backend/internal/tasks"
	"github.com/shmaloogles/business-task-platform/backend/internal/teams"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://shmaloogles:shmaloogles@localhost:5432/shmaloogles?sslmode=disable"
	}

	connectContext, cancelConnect := context.WithTimeout(context.Background(), 5*time.Second)
	db, err := database.Open(connectContext, databaseURL)
	cancelConnect()
	if err != nil {
		logger.Error("could not connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           server.New(db, tasks.NewStore(db), teams.NewStore(db)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("API server started", "address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("API server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdownSignal.Done()
	logger.Info("shutting down API server")

	shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		logger.Error("could not gracefully shut down API server", "error", err)
		os.Exit(1)
	}
}

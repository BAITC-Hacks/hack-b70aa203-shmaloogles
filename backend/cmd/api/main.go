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

	"github.com/shmaloogles/business-task-platform/backend/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           server.New(),
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

package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	loadEnv()

	app, cleanup, err := NewApp()
	if err != nil {
		slog.Error("error while initializing app", slog.Any("error", err))
		os.Exit(1)
	}

	if cleanup != nil {
		defer cleanup()
	}

	slog.SetDefault(app.logger)

	var port string
	if app.config.App.Env == "production" {
		port = os.Getenv("PORT")
	} else {
		port = app.config.App.Port
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.router,
	}

	go func() {
		slog.Info("starting HTTP server", slog.String("port", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed to start", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.Any("error", err))
	}

	slog.Info("server exited gracefully")
}

func loadEnv() {
	slog.Info("loading .env file...")
	file, err := os.Open(".env")
	if err != nil {
		slog.Warn("failed to load .env, relying on system environment variables")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// skip empty and commented line
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// split key value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("error while reading .env", slog.Any("error", err))
		os.Exit(1)
	}
}

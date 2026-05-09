package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/database"
	"server/handlers"
	"server/middleware"
	"server/repository"
	"server/routes"
)

func main() {
	// Logger setup
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Environment vars
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "note_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	serverPort := getEnv("SERVER_PORT", "8080")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, dbSSLMode)

	// Retry koneksi database dengan Context timeout
	ctx, cancelInit := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelInit()

	db, err := database.Connect(ctx, connStr)
	if err != nil {
		slog.Error("Gagal terkoneksi ke database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Database terkoneksi & migrasi berhasil")

	// Setup lapisan layer
	repo := &repository.NoteRepository{DB: db}
	handler := &handlers.NoteHandler{Repo: repo}
	router := routes.SetupRoutes(handler)

	// Terapkan middleware
	wrappedRouter := middleware.LoggingAndRecovery(router)

	// HTTP Server
	srv := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      wrappedRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Graceful Shutdown
	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Server berjalan", "port", serverPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	<-stopCtx.Done()
	slog.Info("Mematikan server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error saat shutdown", "error", err)
	}
	slog.Info("Server berhasil dihentikan")
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

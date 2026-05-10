package main

import (
	"log/slog"
	"net/http"
	"os"
	"server/database"
	"server/handlers"
	"server/repository"
	"server/routes"
)

func main() {
	// Logger terbaru (slog)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 1. Inisialisasi database SQLite
	db, err := database.InitDB("./data/tasks.db")
	if err != nil {
		logger.Error("Gagal inisialisasi database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// 2. Setup repository dan handler
	repo := &repository.TaskRepository{DB: db}
	handler := &handlers.TaskHandler{Repo: repo}

	// 3. Routing menggunakan method pattern
	mux := routes.SetupRoutes(handler)

	// 4. Jalankan server
	port := ":8080"
	logger.Info("Server berjalan", "port", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		logger.Error("Server berhenti", "error", err)
	}
}

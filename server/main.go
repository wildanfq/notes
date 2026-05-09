package main

import (
	"log"
	"net/http"
	"server/handlers"
	"server/middleware"
	"server/storage"
	"time"
)

func main() {
	// init storage
	store := storage.NewMemoryStore()

	// init handler
	noteHandler := handlers.NewNoteHandler(store)

	// router dengan pattern
	mux := http.NewServeMux()
	mux.HandleFunc("GET /notes", noteHandler.GetNotes)
	mux.HandleFunc("POST /notes", noteHandler.CreateNote)
	mux.HandleFunc("PUT /notes/{id}", noteHandler.UpdateNote)
	mux.HandleFunc("DELETE /notes/{id}", noteHandler.DeleteNote)

	// middleware CORS
	handler := middleware.CORS(mux)

	// server config
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server berjalan di http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server gagal: %v", err)
	}
}

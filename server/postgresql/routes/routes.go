package routes

import (
	"net/http"

	"server/handlers"
)

func SetupRoutes(h *handlers.NoteHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Notes endpoints
	mux.HandleFunc("GET /notes", h.List)
	mux.HandleFunc("GET /notes/{id}", h.Get)
	mux.HandleFunc("POST /notes", h.Create)
	mux.HandleFunc("PUT /notes/{id}", h.Update)
	mux.HandleFunc("DELETE /notes/{id}", h.Delete)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}

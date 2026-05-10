package routes

import (
	"net/http"
	"server/handlers"
)

// SetupRoutes mendaftarkan semua endpoint dengan pattern method
func SetupRoutes(h *handlers.TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /notes", h.List)
	mux.HandleFunc("GET /notes/{id}", h.Get)
	mux.HandleFunc("POST /notes", h.Create)
	mux.HandleFunc("PUT /notes/{id}", h.Update)
	mux.HandleFunc("DELETE /notes/{id}", h.Delete)

	return mux
}

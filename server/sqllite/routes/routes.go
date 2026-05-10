package routes

import (
	"net/http"
	"server/handlers"
)

// SetupRoutes mendaftarkan semua endpoint dengan pattern method
func SetupRoutes(h *handlers.TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", h.List)
	mux.HandleFunc("GET /tasks/{id}", h.Get)
	mux.HandleFunc("POST /tasks", h.Create)
	mux.HandleFunc("PUT /tasks/{id}", h.Update)
	mux.HandleFunc("DELETE /tasks/{id}", h.Delete)

	return mux
}

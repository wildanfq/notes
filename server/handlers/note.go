package handlers

import (
	"encoding/json"
	"net/http"
	"server/models"
	"server/storage"
)

type NoteHandler struct {
	store *storage.MemoryStore
}

func NewNoteHandler(store *storage.MemoryStore) *NoteHandler {
	return &NoteHandler{store: store}
}

func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	notes := h.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var newNote models.Note
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if newNote.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	created := h.store.Create(newNote)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updated models.Note
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if updated.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	result, err := h.store.Update(id, updated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

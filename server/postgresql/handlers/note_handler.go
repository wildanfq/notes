package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"server/models"
	"server/repository"
)

type NoteHandler struct {
	Repo *repository.NoteRepository
}

type noteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	notes, err := h.Repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Agar json yang kosong direturn sebagai [] bukan null
	if notes == nil {
		notes = []models.Note{}
	}
	respondJSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	note, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Terjadi kesalahan server", http.StatusInternalServerError)
		return
	}
	if note == nil {
		http.Error(w, "Note tidak ditemukan", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format request tidak valid", http.StatusBadRequest)
		return
	}
	if req.Title == "" || req.Content == "" {
		http.Error(w, "Title dan Content wajib diisi", http.StatusBadRequest)
		return
	}

	note := &models.Note{Title: req.Title, Content: req.Content}
	if err := h.Repo.Create(r.Context(), note); err != nil {
		http.Error(w, "Gagal menyimpan note", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusCreated, note)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format request tidak valid", http.StatusBadRequest)
		return
	}

	note := &models.Note{ID: id, Title: req.Title, Content: req.Content}
	err = h.Repo.Update(r.Context(), note)

	if err == sql.ErrNoRows {
		http.Error(w, "Note tidak ditemukan", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Gagal mengupdate note", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Note tidak ditemukan", http.StatusNotFound)
			return
		}
		http.Error(w, "Gagal menghapus note", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

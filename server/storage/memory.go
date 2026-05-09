package storage

import (
	"errors"
	"server/models"
	"strconv"
	"sync"
)

type MemoryStore struct {
	mu    sync.RWMutex
	notes []models.Note
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		notes:  []models.Note{},
		nextID: 1,
	}
}

func (s *MemoryStore) GetAll() []models.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// return copy
	copyNotes := make([]models.Note, len(s.notes))
	copy(copyNotes, s.notes)
	return copyNotes
}

func (s *MemoryStore) Create(note models.Note) models.Note {
	s.mu.Lock()
	defer s.mu.Unlock()
	note.ID = strconv.Itoa(s.nextID)
	s.nextID++
	s.notes = append(s.notes, note)
	return note
}

func (s *MemoryStore) Update(id string, updated models.Note) (models.Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, n := range s.notes {
		if n.ID == id {
			s.notes[i].Title = updated.Title
			s.notes[i].Content = updated.Content
			return s.notes[i], nil
		}
	}
	return models.Note{}, errors.New("note not found")
}

func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, n := range s.notes {
		if n.ID == id {
			s.notes = append(s.notes[:i], s.notes[i+1:]...)
			return nil
		}
	}
	return errors.New("note not found")
}

package repository

import (
	"context"
	"database/sql"

	"server/models"
)

type NoteRepository struct {
	DB *sql.DB
}

func (r *NoteRepository) GetAll(ctx context.Context) ([]models.Note, error) {
	query := "SELECT id, title, content, created_at, updated_at FROM notes ORDER BY id"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *NoteRepository) GetByID(ctx context.Context, id int) (*models.Note, error) {
	var n models.Note
	query := "SELECT id, title, content, created_at, updated_at FROM notes WHERE id = $1"
	err := r.DB.QueryRowContext(ctx, query, id).
		Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Return nil jika tidak ditemukan
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NoteRepository) Create(ctx context.Context, n *models.Note) error {
	query := `INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.DB.QueryRowContext(ctx, query, n.Title, n.Content).
		Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
}

func (r *NoteRepository) Update(ctx context.Context, n *models.Note) error {
	query := `UPDATE notes SET title = $1, content = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at`
	err := r.DB.QueryRowContext(ctx, query, n.Title, n.Content, n.ID).Scan(&n.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil // Akan ditangani oleh handler sebagai Not Found
	}
	return err
}

func (r *NoteRepository) Delete(ctx context.Context, id int) error {
	res, err := r.DB.ExecContext(ctx, "DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

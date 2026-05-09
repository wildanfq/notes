package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect membuka koneksi ke PostgreSQL dan memastikan tabel tersedia
func Connect(ctx context.Context, connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("gagal buka koneksi: %w", err)
	}

	// Konfigurasi connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test koneksi
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database tidak merespon: %w", err)
	}

	// Buat tabel notes jika belum ada
	schema := `
	CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_notes_title ON notes(title);
	`
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, fmt.Errorf("gagal migrasi tabel: %w", err)
	}

	return db, nil
}

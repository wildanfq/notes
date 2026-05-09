package repository

import (
	"database/sql"
	"server/models"
)

type TaskRepository struct {
	DB *sql.DB
}

// GetAll mengambil semua task
func (r *TaskRepository) GetAll() ([]models.Task, error) {
	rows, err := r.DB.Query("SELECT id, title, done FROM tasks ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// GetByID mengambil satu task berdasarkan ID
func (r *TaskRepository) GetByID(id int) (*models.Task, error) {
	var t models.Task
	row := r.DB.QueryRow("SELECT id, title, done FROM tasks WHERE id = ?", id)
	err := row.Scan(&t.ID, &t.Title, &t.Done)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Create menambah task baru
func (r *TaskRepository) Create(t *models.Task) error {
	res, err := r.DB.Exec("INSERT INTO tasks (title, done) VALUES (?, ?)", t.Title, t.Done)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = int(id)
	return nil
}

// Update mengganti seluruh data task
func (r *TaskRepository) Update(t *models.Task) error {
	_, err := r.DB.Exec("UPDATE tasks SET title = ?, done = ? WHERE id = ?", t.Title, t.Done, t.ID)
	return err
}

// Delete menghapus task berdasarkan ID
func (r *TaskRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

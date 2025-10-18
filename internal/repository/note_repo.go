package repository

import (
	"database/sql"
	"github.com/OscarMarulanda/commentsClean/internal/domain"
)

type NoteRepository interface {
	GetAllByUser(userID int) ([]domain.Note, error)
	Create(note *domain.Note) error
	Update(note *domain.Note) error
	Delete(id int) error
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) GetAllByUser(userID int) ([]domain.Note, error) {
	rows, err := r.db.Query(`SELECT id, title, content, user_id, created_at, updated_at, deleted_at FROM notes WHERE user_id=$1 AND deleted_at IS NULL`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []domain.Note
	for rows.Next() {
		var n domain.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.UserID, &n.CreatedAt, &n.UpdatedAt, &n.DeletedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *noteRepository) Create(note *domain.Note) error {
	query := `
		INSERT INTO notes (title, content, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, note.Title, note.Content, note.UserID).
		Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
}

func (r *noteRepository) Update(note *domain.Note) error {
	query := `
		UPDATE notes SET title=$1, content=$2, updated_at=NOW()
		WHERE id=$3 RETURNING updated_at
	`
	return r.db.QueryRow(query, note.Title, note.Content, note.ID).Scan(&note.UpdatedAt)
}

func (r *noteRepository) Delete(id int) error {
	query := `UPDATE notes SET deleted_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(query, id)
	return err
}
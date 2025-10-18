package repository

import (
	"database/sql"
	"github.com/OscarMarulanda/commentsClean/internal/domain"
)

type UserRepository interface {
	GetByID(id int) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE id=$1`
	row := r.db.QueryRow(query, id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE email=$1`
	row := r.db.QueryRow(query, email)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, created_at)
		VALUES ($1, $2, $3, NOW()) RETURNING id, created_at
	`
	return r.db.QueryRow(query, user.Name, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
}
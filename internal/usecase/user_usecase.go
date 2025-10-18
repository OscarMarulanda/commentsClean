package usecase

import (
	"errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/OscarMarulanda/commentsClean/internal/domain"
	"github.com/OscarMarulanda/commentsClean/internal/repository"
)

type UserUseCase interface {
	Register(name, email, password string) (*domain.User, error)
	GetByID(id int) (*domain.User, error)
	Login(email, password string) (*domain.User, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) UserUseCase {
	return &userUseCase{repo: r}
}

func (u *userUseCase) Register(name, email, password string) (*domain.User, error) {
	existing, _ := u.repo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) GetByID(id int) (*domain.User, error) {
	return u.repo.GetByID(id)
}

// NOTE: in a real app, you'd verify password hash here
func (u *userUseCase) Login(email, password string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
package usecase

import (
	"errors"

	"github.com/OscarMarulanda/commentsClean/internal/domain"
	"github.com/OscarMarulanda/commentsClean/internal/repository"
)

type NoteUseCase interface {
	GetNotesByUser(userID int) ([]domain.Note, error)
	Create(note *domain.Note) error
	Update(note *domain.Note) error
	Delete(noteID int, userID int) error
	SearchNotes(userID int, keyword string) ([]domain.Note, error)
}

type noteUseCase struct {
	repo repository.NoteRepository
}

func NewNoteUseCase(r repository.NoteRepository) NoteUseCase {
	return &noteUseCase{repo: r}
}

func (n *noteUseCase) GetNotesByUser(userID int) ([]domain.Note, error) {
	return n.repo.GetAllByUser(userID)
}

func (n *noteUseCase) SearchNotes(userID int, keyword string) ([]domain.Note, error) {
	if keyword == "" {
		return nil, errors.New("keyword is required")
	}
	return n.repo.SearchByKeyword(userID, keyword)
}

func (n *noteUseCase) Create(note *domain.Note) error {
	if note.Title == "" {
		return errors.New("title is required")
	}
	return n.repo.Create(note)
}

func (n *noteUseCase) Update(note *domain.Note) error {
	if note.ID == 0 {
		return errors.New("note ID required")
	}
	return n.repo.Update(note)
}

func (n *noteUseCase) Delete(noteID int, userID int) error {
	return n.repo.Delete(noteID)
}
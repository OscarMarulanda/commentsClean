package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/OscarMarulanda/commentsClean/internal/domain"
	"github.com/OscarMarulanda/commentsClean/internal/usecase"
	"github.com/OscarMarulanda/commentsClean/internal/middleware"
)

type NoteHandler struct {
	usecase usecase.NoteUseCase
}

func NewNoteHandler(u usecase.NoteUseCase) *NoteHandler {
	return &NoteHandler{usecase: u}
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.ExtractUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notes, err := h.usecase.GetNotesByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) Search(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.ExtractUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Error(w, "keyword query param required", http.StatusBadRequest)
		return
	}

	notes, err := h.usecase.SearchNotes(userID, keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.ExtractUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var note domain.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note.UserID = userID

	if err := h.usecase.Create(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	var note domain.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.usecase.Update(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := h.usecase.Delete(id, 0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
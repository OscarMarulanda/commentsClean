package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/OscarMarulanda/commentsClean/internal/handler"
	"github.com/OscarMarulanda/commentsClean/internal/middleware"
)

func SetupRouter(userHandler *handler.UserHandler, noteHandler *handler.NoteHandler) http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	// User routes
	api.HandleFunc("/users/register", userHandler.Register).Methods("POST")
	api.HandleFunc("/users/login", userHandler.Login).Methods("POST")

	// Note routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.JWTAuth)
	protected.HandleFunc("/notes/search", noteHandler.Search).Methods("GET")
	protected.HandleFunc("/notes", noteHandler.GetAll).Methods("GET")
	protected.HandleFunc("/notes", noteHandler.Create).Methods("POST")
	protected.HandleFunc("/notes", noteHandler.Update).Methods("PUT")
	protected.HandleFunc("/notes/{id}", noteHandler.Delete).Methods("DELETE")

	return r
}
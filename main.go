package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/OscarMarulanda/commentsClean/internal/db"
	"github.com/OscarMarulanda/commentsClean/internal/handler"
	"github.com/OscarMarulanda/commentsClean/internal/repository"
	"github.com/OscarMarulanda/commentsClean/internal/router"
	"github.com/OscarMarulanda/commentsClean/internal/usecase"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found")
	}

	// Connect to database
	conn, err := db.Connect()
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}
	defer conn.Close()

	// Initialize layers
	userRepo := repository.NewUserRepository(conn)
	noteRepo := repository.NewNoteRepository(conn)

	userUC := usecase.NewUserUseCase(userRepo)
	noteUC := usecase.NewNoteUseCase(noteRepo)

	userHandler := handler.NewUserHandler(userUC)
	noteHandler := handler.NewNoteHandler(noteUC)

	// Setup routes
	r := router.SetupRouter(userHandler, noteHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on :%s\n", port)
	http.ListenAndServe(":"+port, r)
}
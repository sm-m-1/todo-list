package main

import (
	"log"
	"net/http"
	"time"
	"todo-list/config"
	"todo-list/internal/handlers"
	"todo-list/internal/repos"
	"todo-list/internal/services"
	"todo-list/pkg/database"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize database
	db := database.InitDB()
	defer database.CloseDB(db)

	// Initialize session manager
	sessionManager := scs.New()
	sessionManager.Store = database.NewGORMStore(db, 24*time.Hour)
	sessionManager.Lifetime = 24 * time.Hour

	sessionManager.Cookie.Name = "session_token"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = false // Set to true in production for HTTPS
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	todoRepo := repos.NewTodoRepository(db)
	userRepo := repos.NewUserRepository(db)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)
	authHandler := handlers.NewAuthHandler(userRepo, sessionManager)

	// Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(sessionManager.LoadAndSave)

	// Routes
	r.Post("/register", authHandler.Register())
	r.Post("/login", authHandler.Login())
	r.Post("/logout", authHandler.Logout())
	r.Get("/home", config.SessionMiddleware(handlers.Home(), sessionManager))
	r.Get("/todos", config.SessionMiddleware(todoHandler.GetTodos(), sessionManager))
	r.Post("/todos", config.SessionMiddleware(todoHandler.CreateTodo(), sessionManager))
	r.Put("/todos/{id}", config.SessionMiddleware(todoHandler.UpdateTodo(), sessionManager))
	r.Delete("/todos/{id}", config.SessionMiddleware(todoHandler.DeleteTodo(), sessionManager))
	r.Get("/todos/{id}", config.SessionMiddleware(todoHandler.GetTodo(), sessionManager))

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServer(":8080", sessionManager.LoadAndSave(r)))
	// log.Fatal(http.ListenAndServe(":8080", r))
}

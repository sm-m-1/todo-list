package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"todo-list/database"
	"todo-list/handlers"
	"todo-list/repositories"
	"todo-list/services"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize database
	db := database.InitDB()
	defer database.CloseDB(db)

	todoRepo := repositories.NewTodoRepository(db)
	todoService := services.NewTodoService(todoRepo)

	// Initialize session manager
	sessionManager := scs.New()
	sessionManager.Store = database.NewGORMStore(db, 24*time.Hour)
	sessionManager.Lifetime = 24 * time.Hour

	sessionManager.Cookie.Name = "session_token"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = false // Set to true in production for HTTPS
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	// Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(sessionManager.LoadAndSave)

	// Routes
	r.Post("/register", handlers.Register(db))
	r.Get("/home", SessionMiddleware(handlers.Home(), sessionManager))
	r.Post("/login", handlers.Login(db, sessionManager))
	r.Post("/logout", handlers.Logout(sessionManager))
	r.Get("/todos", SessionMiddleware(handlers.GetTodos(todoService), sessionManager))
	r.Post("/todos", SessionMiddleware(handlers.CreateTodo(todoService), sessionManager))
	r.Put("/todos/{id}", SessionMiddleware(handlers.UpdateTodo(todoService), sessionManager))
	r.Delete("/todos/{id}", SessionMiddleware(handlers.DeleteTodo(todoService), sessionManager))

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", sessionManager.LoadAndSave(r)))
	// log.Fatal(http.ListenAndServe(":8080", r))
}

// SessionMiddleware ensures the user is authenticated
func SessionMiddleware(next http.HandlerFunc, sessionManager *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		session := sessionManager.GetString(r.Context(), "userID")
		fmt.Println("session value from db after in middleware:::: ", session)

		if session == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

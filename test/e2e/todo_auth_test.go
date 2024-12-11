package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"todo-list/config"
	"todo-list/internal/handlers"
	"todo-list/internal/models"
	"todo-list/internal/repos"
	"todo-list/internal/services"
	"todo-list/pkg/database"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTodoAppWithAuth(t *testing.T) {
	// Set up the test database and router
	db, err := setupTestDatabase() // Use SQLite in-memory or mock DB
	if err != nil {
		t.Fatalf("Failed to set up database: %v", err)
	}
	// defer db.Close()

	router := setupRouter(db) // Set up your chi router with all handlers

	// Start the test server
	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{}
	var authCookie *http.Cookie // Cookie for session tracking

	// Step 1: Test User Registration
	registerPayload := map[string]interface{}{
		"username": "testuser",
		"password": "password123",
	}
	registerBody, _ := json.Marshal(registerPayload)

	registerReq, _ := http.NewRequest("POST", server.URL+"/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	registerResp, err := client.Do(registerReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, registerResp.StatusCode)

	// Step 2: Test User Login
	loginPayload := map[string]interface{}{
		"username": "testuser",
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginPayload)

	loginReq, _ := http.NewRequest("POST", server.URL+"/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err := client.Do(loginReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	// Extract session cookie
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "session" {
			authCookie = cookie
		}
	}
	assert.NotNil(t, authCookie)

	// Step 3: Test Creating a Todo (Authenticated)
	createPayload := map[string]interface{}{
		"title":       "Test Todo",
		"description": "This is a test todo",
	}
	createBody, _ := json.Marshal(createPayload)

	createReq, _ := http.NewRequest("POST", server.URL+"/todos", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.AddCookie(authCookie) // Attach session cookie

	createResp, err := client.Do(createReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, createResp.StatusCode)

	var createdTodo models.Todo
	json.NewDecoder(createResp.Body).Decode(&createdTodo)
	createResp.Body.Close()

	fmt.Println("createdTodo: ", createdTodo)
	fmt.Println("")

	assert.Equal(t, "Test Todo", createdTodo.Title)

	// Step 4: Test Logout
	logoutReq, _ := http.NewRequest("POST", server.URL+"/logout", nil)
	logoutReq.AddCookie(authCookie) // Attach session cookie

	logoutResp, err := client.Do(logoutReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, logoutResp.StatusCode)

	// Step 5: Test Unauthorized Access After Logout
	todoId := strconv.FormatInt(int64(createdTodo.ID), 10)
	getReq, _ := http.NewRequest("GET", server.URL+"/todos/"+todoId, nil)
	getReq.AddCookie(authCookie) // Use the same session cookie

	getResp, err := client.Do(getReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, getResp.StatusCode)
}

func setupTestDatabase() (*gorm.DB, error) {
	// Use SQLite in-memory for testing or a mock DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate schemas
	db.AutoMigrate(&models.User{}, &models.Todo{}, &models.Session{})
	return db, nil
}

func setupRouter(db *gorm.DB) http.Handler {
	// Set up your chi router and handlers
	r := chi.NewRouter()

	// Middleware for session management
	sessionManager := scs.New()
	sessionManager.Store = database.NewGORMStore(db, 24*time.Hour)

	r.Use(sessionManager.LoadAndSave)

	// Initialize handlers
	userRepo := repos.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo, sessionManager)

	todoRepo := repos.NewTodoRepository(db)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	// User routes
	r.Post("/register", authHandler.Register())
	r.Post("/login", authHandler.Login())
	r.Post("/logout", authHandler.Logout())

	// Todo routes
	r.Group(func(r chi.Router) {
		// r.Use(config.SessionMiddleware(sessionManager)) // Protect routes with auth middleware
		r.Post("/todos", config.SessionMiddleware(todoHandler.CreateTodo(), sessionManager))
		r.Put("/todos/{id}", config.SessionMiddleware(todoHandler.UpdateTodo(), sessionManager))
		r.Delete("/todos/{id}", config.SessionMiddleware(todoHandler.DeleteTodo(), sessionManager))
		r.Get("/todos/{id}", config.SessionMiddleware(todoHandler.GetTodo(), sessionManager))
	})

	return r
}

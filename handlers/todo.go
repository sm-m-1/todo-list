package handlers

import (
	"encoding/json"
	"net/http"
	"todo-list/models"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// GetTodos retrieves all todos for the authenticated user
func GetTodos(db *gorm.DB, sessionManager *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var todos []models.Todo
		if err := db.Find(&todos).Error; err != nil {
			http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	}
}

// CreateTodo adds a new todo
func CreateTodo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo models.Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := db.Create(&todo).Error; err != nil {
			http.Error(w, "Failed to create todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	}
}

// UpdateTodo updates an existing todo
func UpdateTodo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var todo models.Todo
		if err := db.First(&todo, id).Error; err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := db.Save(&todo).Error; err != nil {
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(todo)
	}
}

// DeleteTodo removes a todo
func DeleteTodo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := db.Delete(&models.Todo{}, id).Error; err != nil {
			http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

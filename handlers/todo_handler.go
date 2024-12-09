package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo-list/models"
	"todo-list/services"

	"github.com/go-chi/chi/v5"
)

// GetTodos retrieves all todos for the authenticated user
func GetTodos(service *services.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		todos, err := service.GetTodoList()

		if err != nil {
			http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	}
}

// CreateTodo adds a new todo
func CreateTodo(service *services.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(uint)
		fmt.Println("context userID in CreateTodo:::: ", userID, ok)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var todo models.Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Associate todo with logged-in user
		todo.UserID = userID

		if err := service.AddTodo(&todo); err != nil {
			fmt.Println("error when trying to add todo ", todo, err)
			http.Error(w, "Failed to create todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	}
}

// UpdateTodo updates an existing todo
func UpdateTodo(service *services.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		todo, err := service.GetTodo(id)
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := service.EditTodo(&todo); err != nil {
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(todo)
	}
}

// DeleteTodo removes a todo
func DeleteTodo(service *services.TodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := service.RemoveTodo(id); err != nil {
			http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

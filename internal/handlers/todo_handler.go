package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo-list/internal/models"
	"todo-list/internal/services"

	"github.com/go-chi/chi/v5"
)

type TodoHandler struct {
	service *services.TodoService
}

func NewTodoHandler(service *services.TodoService) *TodoHandler {
	return &TodoHandler{service}
}

// GetTodos retrieves all todos for the authenticated user
func (h *TodoHandler) GetTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(uint)
		// missing userID in the request context, which should exist from being set in SessionMiddleware
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		todos, err := h.service.GetTodoList(userID)

		if err != nil {
			http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	}
}

// CreateTodo adds a new todo
func (h *TodoHandler) CreateTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(uint)
		// missing userID in the request context, which should exist from being set in SessionMiddleware
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

		if err := h.service.AddTodo(&todo); err != nil {
			fmt.Println("error when trying to add todo ", todo, err)
			http.Error(w, "Failed to create todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	}
}

// UpdateTodo updates an existing todo
func (h *TodoHandler) UpdateTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		todo, err := h.service.GetTodo(id)
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		userID, ok := r.Context().Value("userID").(uint)
		// missing userID in the request context, which should exist from being set in SessionMiddleware

		if !ok || todo.UserID != userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := h.service.EditTodo(&todo); err != nil {
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(todo)
	}
}

// DeleteTodo removes a todo
func (h *TodoHandler) DeleteTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		todo, err := h.service.GetTodo(id)
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		userID, ok := r.Context().Value("userID").(uint)
		// missing userID in the request context, which should exist from being set in SessionMiddleware

		if !ok || todo.UserID != userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := h.service.RemoveTodo(id); err != nil {
			http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *TodoHandler) GetTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		todo, err := h.service.GetTodo(id)
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(todo)
	}
}

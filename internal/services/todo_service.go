package services

import (
	"todo-list/internal/models"
	"todo-list/internal/repos"
)

type TodoService struct {
	repo *repos.TodoRepository
}

// the constructor for TodoService

func NewTodoService(repo *repos.TodoRepository) *TodoService {
	return &TodoService{repo}
}

func (s *TodoService) GetTodoList(userId uint) ([]models.Todo, error) {
	return s.repo.GetAllTodos(userId)
}

func (s *TodoService) AddTodo(todo *models.Todo) error {
	return s.repo.CreateTodo(todo)
}

func (s *TodoService) EditTodo(todo *models.Todo) error {
	return s.repo.UpdateTodo(todo)
}

func (s *TodoService) RemoveTodo(id string) error {
	return s.repo.DeleteTodo(id)
}

func (s *TodoService) GetTodo(id string) (models.Todo, error) {
	return s.repo.GetTodo(id)
}

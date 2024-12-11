package repos

import (
	"todo-list/internal/models"

	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

// Constructor for TodoRepository
func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db}
}

// Fetch all todos
func (r *TodoRepository) GetAllTodos(userId uint) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Where("user_id = ?", userId).Find(&todos).Error
	return todos, err
}

// Save a new todo
func (r *TodoRepository) CreateTodo(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

// update a todo
func (r *TodoRepository) UpdateTodo(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

// delete a todo
func (r *TodoRepository) DeleteTodo(id string) error {
	return r.db.Delete(&models.Todo{}, id).Error
}

// get a todo
func (r *TodoRepository) GetTodo(id string) (models.Todo, error) {
	var todo models.Todo
	return todo, r.db.First(&todo, id).Error
}

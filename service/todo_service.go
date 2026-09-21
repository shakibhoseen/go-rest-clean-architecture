package service

import (
	"context"
	"crud/models"
	"crud/repository"
	"errors"
)

type TodoService interface {
	CreateTodo(ctx context.Context, todo *models.Todo) error
	GetTodos(ctx context.Context, userID int) ([]models.Todo, error)
	GetTodoByID(ctx context.Context, id, userID int) (*models.Todo, error)
	UpdateTodo(ctx context.Context, todo *models.Todo) error
	DeleteTodo(ctx context.Context, id, userID int) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(ctx context.Context, todo *models.Todo) error {
	return s.repo.Create(ctx, todo)
}

func (s *todoService) GetTodos(ctx context.Context, userID int) ([]models.Todo, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *todoService) GetTodoByID(ctx context.Context, id, userID int) (*models.Todo, error) {
	todo, err := s.repo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

func (s *todoService) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	return s.repo.Update(ctx, todo)
}

func (s *todoService) DeleteTodo(ctx context.Context, id, userID int) error {
	return s.repo.Delete(ctx, id, userID)
}

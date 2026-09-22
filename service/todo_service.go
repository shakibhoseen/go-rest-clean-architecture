package service

import (
	"context"
	"crud/models"
	"crud/repository"
	"crud/utils"
	"errors"
)

type TodoService interface {
	CreateTodo(ctx context.Context, todo *models.Todo) error
	GetTodos(ctx context.Context, userID int, completed *bool, sortOrder string, page, limit int) (*utils.PaginatedResponse, error)
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

func (s *todoService) GetTodos(ctx context.Context, userID int, completed *bool, sortOrder string, page, limit int) (*utils.PaginatedResponse, error) {
	// ডিফল্ট মান হ্যান্ডলিং
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 { // ক্লায়েন্ট যাতে একবারে ১০,০০০ না চেয়ে বসে
		limit = 10
	}

	// অফসেট ক্যালকুলেশন
	offset := (page - 1) * limit

	todos, totalRecords, err := s.repo.GetByUserID(ctx, userID, completed, sortOrder, limit, offset)
	if err != nil {
		return nil, err
	}

	// Total Pages হিসেব: ceil(totalRecords / limit)
	totalPages := (totalRecords + limit - 1) / limit

	return &utils.PaginatedResponse{
		TotalRecords: totalRecords,
		CurrentPage:  page,
		TotalPages:   totalPages,
		Limit:        limit,
		Data:         todos,
	}, nil
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

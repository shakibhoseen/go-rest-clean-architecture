package service

import (
	"context"
	"crud/models"
	"crud/repository"
	"errors"
)

type UserService interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, u *models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, u *models.User) error {
	return s.repo.Create(ctx, u)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (s *userService) UpdateUser(ctx context.Context, u *models.User) error {
	return s.repo.Update(ctx, u)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

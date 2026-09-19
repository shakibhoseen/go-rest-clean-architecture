package service

import (
	"context"
	"crud/models"
	"crud/repository"
	"crud/utils"
	"errors"
)

type UserService interface {
	RegisterUser(ctx context.Context, u *models.User) error
	Login(ctx context.Context, email, password string) (string, *models.User, error)
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

func (s *userService) RegisterUser(ctx context.Context, u *models.User) error {
	// 1. Password hash kora
	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	u.Password = "" // memory theke plain text clear kora

	return s.repo.Create(ctx, u)
}

func (s *userService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Password match check
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("invalid email or password")
	}

	// JWT token generate
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}

	user.PasswordHash = "" // response-e jeno hash na jay
	return token, user, nil
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

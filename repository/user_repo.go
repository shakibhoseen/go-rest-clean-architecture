package repository

import (
	"context"
	"crud/models"
	"database/sql"
	"errors"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id int) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, u.Name, u.Email, u.PasswordHash).Scan(&u.ID)
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, name, email, password_hash FROM users WHERE email = $1`
	var u models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, err
}

// GetAllUsers handles GET /users
func (r *userRepo) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, email FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`
	var u models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // ইউজার না পেলে nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	res, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

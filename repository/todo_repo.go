package repository

import (
	"context"
	"crud/models"
	"database/sql"
	"errors"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *models.Todo) error
	GetByUserID(ctx context.Context, userID int, limit, offset int) ([]models.Todo, int, error)
	GetByIDAndUserID(ctx context.Context, id, userID int) (*models.Todo, error)
	Update(ctx context.Context, todo *models.Todo) error
	Delete(ctx context.Context, id, userID int) error
}

type todoRepo struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepo{db: db}
}

func (r *todoRepo) Create(ctx context.Context, todo *models.Todo) error {
	query := `INSERT INTO todos (user_id, title, completed) 
	          VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, todo.UserID, todo.Title, todo.Completed).
		Scan(&todo.ID, &todo.CreatedAt)
}

func (r *todoRepo) GetByUserID(ctx context.Context, userID int, limit, offset int) ([]models.Todo, int, error) {
	// ১. এই ইউজারের মোট টোডো কয়টি আছে তা গণনা
	var totalRecords int
	countQuery := `SELECT COUNT(*) FROM todos WHERE user_id = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalRecords); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, user_id, title, completed, created_at 
	          FROM todos WHERE user_id = $1 ORDER BY created_at DESC
			  LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return todos, totalRecords, nil
}

func (r *todoRepo) GetByIDAndUserID(ctx context.Context, id, userID int) (*models.Todo, error) {
	query := `SELECT id, user_id, title, completed, created_at 
	          FROM todos WHERE id = $1 AND user_id = $2`
	var t models.Todo
	err := r.db.QueryRowContext(ctx, query, id, userID).
		Scan(&t.ID, &t.UserID, &t.Title, &t.Completed, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *todoRepo) Update(ctx context.Context, todo *models.Todo) error {
	query := `UPDATE todos SET title = $1, completed = $2 
	          WHERE id = $3 AND user_id = $4`
	res, err := r.db.ExecContext(ctx, query, todo.Title, todo.Completed, todo.ID, todo.UserID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("todo not found or unauthorized")
	}
	return nil
}

func (r *todoRepo) Delete(ctx context.Context, id, userID int) error {
	query := `DELETE FROM todos WHERE id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("todo not found or unauthorized")
	}
	return nil
}

package models

import (
	"errors"
	"strings"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Todo) Validate() error {
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		return errors.New("title is required")
	}
	if len(t.Title) < 3 {
		return errors.New("title must be at least 3 characters")
	}
	return nil
}

package models

import (
	"errors"
	"net/mail"
	"strings"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) Validate() error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)

	if u.Name == "" {
		return errors.New("name is required")
	}
	if len(u.Name) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	// ইমেইল ফরম্যাট সঠিক কি না বিল্ট-ইন net/mail প্যাকেজ দিয়ে চেক
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email address format")
	}

	return nil
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

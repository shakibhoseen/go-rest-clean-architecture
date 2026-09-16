package handlers

import (
	"context"
	"crud/models"
	"crud/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// GetAllUsers handles GET /users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	rows, err := h.DB.QueryContext(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {

			utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "Users retrieved successfully", users)
}

// PostUser handles POST /users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {

		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// ১. স্বয়ংক্রিয় ইনপুট ভ্যালিডেশন
	if err := u.Validate(); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	result, err := h.DB.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", u.Name, u.Email)
	if err != nil {

		utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, _ := result.LastInsertId()
	u.ID = int(id)

	utils.JSON(w, http.StatusCreated, "User created successfully", u)
}

// GetUserByID handles GET /users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {

		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var u models.User
	err = h.DB.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.Email)
	if err == sql.ErrNoRows {

		utils.ErrorJSON(w, http.StatusNotFound, "User not found")
		return
	} else if err != nil {

		utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, "User fetch success", u)
}

// DeleteUser handles DELETE /users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {

		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	result, err := h.DB.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {

		utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {

		utils.ErrorJSON(w, http.StatusNotFound, "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// update handlers Put /users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {

		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid user ID")
		//Utils
		return
	}

	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {

		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// আপডেট করার আগেও ভ্যালিডেশন চেক
	if err := u.Validate(); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	result, err := h.DB.ExecContext(ctx, "UPDATE users SET name = ?, email = ? WHERE id = ?", u.Name, u.Email, id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err.Error())

		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	u.ID = id

	utils.JSON(w, http.StatusNoContent, "User Update successfully", u)

}

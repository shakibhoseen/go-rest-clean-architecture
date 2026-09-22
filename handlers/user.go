package handlers

import (
	"context"
	"crud/models"
	"crud/service"
	"crud/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	svc service.UserService
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with name, email, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "User Registration Info"
// @Success      201 {object} utils.APIResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      409 {object} utils.ErrorResponse
// @Router       /auth/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := u.Validate(); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := u.ValidatePassword(); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.svc.RegisterUser(ctx, &u); err != nil {
		utils.ErrorJSON(w, http.StatusConflict, "Email already registered or database error")
		return
	}

	utils.JSON(w, http.StatusCreated, "User registered successfully", u)
}

// Login godoc
// @Summary      User Login
// @Description  Authenticate user with email and password to receive a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Login Credentials"
// @Success      200 {object} utils.APIResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Router       /auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	token, user, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		utils.ErrorJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	respData := map[string]any{
		"token": token,
		"user":  user,
	}

	utils.JSON(w, http.StatusOK, "Login successful", respData)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := u.Validate(); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.svc.CreateUser(ctx, &u); err != nil {
		utils.ErrorJSON(w, http.StatusConflict, "Email already exists or DB error")
		return
	}

	utils.JSON(w, http.StatusCreated, "User created successfully", u)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	users, err := h.svc.GetAllUsers(ctx)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	utils.JSON(w, http.StatusOK, "Users retrieved successfully", users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	u, err := h.svc.GetUserByID(ctx, id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "User fetched successfully", u)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := u.Validate(); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	u.ID = id

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.svc.UpdateUser(ctx, &u); err != nil {
		utils.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "User updated successfully", u)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.svc.DeleteUser(ctx, id); err != nil {
		utils.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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

type TodoHandler struct {
	svc service.TodoService
}

type UpdateTodoRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func NewTodoHandler(svc service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := todo.Validate(); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	todo.UserID = userID

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.svc.CreateTodo(ctx, &todo); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	utils.JSON(w, http.StatusCreated, "Todo created successfully", todo)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	// 1. Context theke authenticated user ID neya
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 2. URL theke todo ID parse kora
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	// 3. Request body decode kora
	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// 4. Todo struct banano & validation
	todo := models.Todo{
		ID:        id,
		UserID:    userID,
		Title:     req.Title,
		Completed: req.Completed,
	}

	if err := todo.Validate(); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// 5. Service layer call kora (jekhane query-te id ebong user_id match kora hoy)
	if err := h.svc.UpdateTodo(ctx, &todo); err != nil {
		utils.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "Todo updated successfully", todo)
}

func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	todos, err := h.svc.GetTodos(ctx, userID)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to retrieve todos")
		return
	}

	utils.JSON(w, http.StatusOK, "Todos retrieved successfully", todos)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.svc.DeleteTodo(ctx, id, userID); err != nil {
		utils.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

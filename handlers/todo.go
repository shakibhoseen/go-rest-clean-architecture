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

// Create godoc
// @Summary      Create a new todo
// @Description  Create a new todo task for the authenticated user
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateTodoRequest true "Todo Payload"
// @Success      201 {object} utils.APIResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Router       /todos [post]
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

// Update godoc
// @Summary      Update a todo
// @Description  Update title or completion status of an existing todo
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path     int               true "Todo ID"
// @Param        request body     UpdateTodoRequest  true "Todo Update Payload"
// @Success      200     {object} utils.APIResponse
// @Failure      400     {object} utils.ErrorResponse
// @Failure      401     {object} utils.ErrorResponse
// @Failure      404     {object} utils.ErrorResponse
// @Router       /todos/{id} [put]
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

// GetAll godoc
// @Summary      Get all todos
// @Description  Get a paginated, filtered, and sorted list of todos for the authenticated user
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Page number" default(1)
// @Param        limit      query     int     false  "Items per page" default(10)
// @Param        completed  query     bool    false  "Filter by completed status"
// @Param        sort       query     string  false  "Sort order (asc/desc)" default(desc)
// @Success      200        {object}  utils.PaginatedResponse
// @Failure      401        {object}  utils.ErrorResponse
// @Failure      500        {object}  utils.ErrorResponse
// @Router       /todos [get]
func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	//url query parameter read default page =1 , limit =10
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	sortOrder := r.URL.Query().Get("sort")

	var completedPtr *bool
	completedStr := r.URL.Query().Get("completed")
	if completedStr != "" {
		if val, err := strconv.ParseBool(completedStr); err == nil {
			completedPtr = &val
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	paginatedData, err := h.svc.GetTodos(ctx, userID, completedPtr, sortOrder, page, limit)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to retrieve todos")
		return
	}

	utils.JSON(w, http.StatusOK, "Todos retrieved successfully", paginatedData)
}

// Delete godoc
// @Summary      Delete a todo
// @Description  Delete a todo belonging to the authenticated user
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Todo ID"
// @Success      204 "No Content"
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Router       /todos/{id} [delete]
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

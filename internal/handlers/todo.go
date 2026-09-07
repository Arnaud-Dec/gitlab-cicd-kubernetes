package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"todoapp/internal/metrics"
	"todoapp/internal/models"
)

type TodoHandler struct {
	DB *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{DB: db}
}

func respondError(c echo.Context, status int, msg string) error {
	return c.JSON(status, errorResponse{Error: msg})
}

func parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// recordCompletionChange adjusts the open/closed gauges when a todo's
// completed state actually changes (a no-op update shouldn't move them).
func recordCompletionChange(wasCompleted, isCompleted bool) {
	if wasCompleted == isCompleted {
		return
	}
	if isCompleted {
		metrics.TodosOpen.Dec()
		metrics.TodosClosed.Inc()
	} else {
		metrics.TodosOpen.Inc()
		metrics.TodosClosed.Dec()
	}
}

// List godoc
//
//	@Summary	List todos
//	@Tags		todos
//	@Produce	json
//	@Success	200	{array}		models.Todo
//	@Failure	500	{object}	errorResponse
//	@Router		/api/todos [get]
func (h *TodoHandler) List(c echo.Context) error {
	var todos []models.Todo
	if err := h.DB.Order("created_at desc").Find(&todos).Error; err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to list todos")
	}
	return c.JSON(http.StatusOK, todos)
}

// Create godoc
//
//	@Summary	Create a todo
//	@Tags		todos
//	@Accept		json
//	@Produce	json
//	@Param		todo	body		CreateTodoRequest	true	"Todo to create"
//	@Success	201		{object}	models.Todo
//	@Failure	400		{object}	errorResponse
//	@Router		/api/todos [post]
func (h *TodoHandler) Create(c echo.Context) error {
	var req CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "title is required")
	}

	todo := models.Todo{Title: req.Title, Description: req.Description}
	if err := h.DB.Create(&todo).Error; err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to create todo")
	}

	metrics.TodosCreatedTotal.Inc()
	metrics.TodosOpen.Inc()

	return c.JSON(http.StatusCreated, todo)
}

// Get godoc
//
//	@Summary	Get a todo by ID
//	@Tags		todos
//	@Produce	json
//	@Param		id	path		int	true	"Todo ID"
//	@Success	200	{object}	models.Todo
//	@Failure	400	{object}	errorResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/api/todos/{id} [get]
func (h *TodoHandler) Get(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid id")
	}

	var todo models.Todo
	if err := h.DB.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return respondError(c, http.StatusNotFound, "todo not found")
		}
		return respondError(c, http.StatusInternalServerError, "failed to get todo")
	}
	return c.JSON(http.StatusOK, todo)
}

// Update godoc
//
//	@Summary	Update a todo
//	@Tags		todos
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int					true	"Todo ID"
//	@Param		todo	body		UpdateTodoRequest	true	"Updated todo fields"
//	@Success	200		{object}	models.Todo
//	@Failure	400		{object}	errorResponse
//	@Failure	404		{object}	errorResponse
//	@Router		/api/todos/{id} [put]
func (h *TodoHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid id")
	}

	var req UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "title is required")
	}

	var todo models.Todo
	if err := h.DB.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return respondError(c, http.StatusNotFound, "todo not found")
		}
		return respondError(c, http.StatusInternalServerError, "failed to get todo")
	}

	wasCompleted := todo.Completed
	todo.Title = req.Title
	todo.Description = req.Description
	todo.Completed = req.Completed
	if err := h.DB.Save(&todo).Error; err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to update todo")
	}

	recordCompletionChange(wasCompleted, todo.Completed)

	return c.JSON(http.StatusOK, todo)
}

// Toggle godoc
//
//	@Summary	Toggle a todo's completed state
//	@Tags		todos
//	@Produce	json
//	@Param		id	path		int	true	"Todo ID"
//	@Success	200	{object}	models.Todo
//	@Failure	400	{object}	errorResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/api/todos/{id}/toggle [patch]
func (h *TodoHandler) Toggle(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid id")
	}

	var todo models.Todo
	if err := h.DB.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return respondError(c, http.StatusNotFound, "todo not found")
		}
		return respondError(c, http.StatusInternalServerError, "failed to get todo")
	}

	wasCompleted := todo.Completed
	todo.Completed = !todo.Completed
	if err := h.DB.Save(&todo).Error; err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to toggle todo")
	}

	recordCompletionChange(wasCompleted, todo.Completed)

	return c.JSON(http.StatusOK, todo)
}

// Delete godoc
//
//	@Summary	Delete a todo
//	@Tags		todos
//	@Param		id	path	int	true	"Todo ID"
//	@Success	204	"No Content"
//	@Failure	400	{object}	errorResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/api/todos/{id} [delete]
func (h *TodoHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid id")
	}

	var todo models.Todo
	if err := h.DB.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return respondError(c, http.StatusNotFound, "todo not found")
		}
		return respondError(c, http.StatusInternalServerError, "failed to get todo")
	}

	if err := h.DB.Delete(&todo).Error; err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to delete todo")
	}

	if todo.Completed {
		metrics.TodosClosed.Dec()
	} else {
		metrics.TodosOpen.Dec()
	}

	return c.NoContent(http.StatusNoContent)
}

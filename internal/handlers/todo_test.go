package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"todoapp/internal/handlers"
	"todoapp/internal/models"
	"todoapp/internal/server"
)

func newEcho() *echo.Echo {
	e := echo.New()
	e.Validator = server.NewValidator()
	return e
}

func doRequest(e *echo.Echo, method, path string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequestWithContext(context.Background(), method, path, reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func createTodo(t *testing.T, e *echo.Echo, title string) models.Todo {
	t.Helper()
	rec := doRequest(e, http.MethodPost, "/api/todos", handlers.CreateTodoRequest{Title: title})
	require.Equal(t, http.StatusCreated, rec.Code)

	var todo models.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &todo))
	return todo
}

func setupRouter(db *gorm.DB) *echo.Echo {
	e := newEcho()
	h := handlers.NewTodoHandler(db)
	e.GET("/api/todos", h.List)
	e.POST("/api/todos", h.Create)
	e.GET("/api/todos/:id", h.Get)
	e.PUT("/api/todos/:id", h.Update)
	e.PATCH("/api/todos/:id/toggle", h.Toggle)
	e.DELETE("/api/todos/:id", h.Delete)
	return e
}

func TestTodoHandler_Create(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)

	rec := doRequest(e, http.MethodPost, "/api/todos", handlers.CreateTodoRequest{Title: "Write plan"})
	assert.Equal(t, http.StatusCreated, rec.Code)

	var todo models.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &todo))
	assert.Equal(t, "Write plan", todo.Title)
	assert.False(t, todo.Completed)
}

func TestTodoHandler_Create_MissingTitle(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)

	rec := doRequest(e, http.MethodPost, "/api/todos", handlers.CreateTodoRequest{Title: ""})
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTodoHandler_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)

	rec := doRequest(e, http.MethodGet, "/api/todos", nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	var todos []models.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &todos))
	assert.Empty(t, todos)
}

func TestTodoHandler_List_Populated(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	createTodo(t, e, "First")
	createTodo(t, e, "Second")

	rec := doRequest(e, http.MethodGet, "/api/todos", nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	var todos []models.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &todos))
	assert.Len(t, todos, 2)
}

func TestTodoHandler_Get_Found(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	created := createTodo(t, e, "Find me")

	rec := doRequest(e, http.MethodGet, "/api/todos/"+itoa(created.ID), nil)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTodoHandler_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)

	rec := doRequest(e, http.MethodGet, "/api/todos/999", nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTodoHandler_Get_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)

	rec := doRequest(e, http.MethodGet, "/api/todos/not-a-number", nil)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTodoHandler_Update(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	created := createTodo(t, e, "Original")

	rec := doRequest(e, http.MethodPut, "/api/todos/"+itoa(created.ID), handlers.UpdateTodoRequest{
		Title:     "Updated",
		Completed: true,
	})
	assert.Equal(t, http.StatusOK, rec.Code)

	var todo models.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &todo))
	assert.Equal(t, "Updated", todo.Title)
	assert.True(t, todo.Completed)
}

func TestTodoHandler_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)

	rec := doRequest(e, http.MethodPut, "/api/todos/999", handlers.UpdateTodoRequest{Title: "Updated"})
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTodoHandler_Toggle(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	created := createTodo(t, e, "Toggle me")
	assert.False(t, created.Completed)

	rec := doRequest(e, http.MethodPatch, "/api/todos/"+itoa(created.ID)+"/toggle", nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	var todo models.Todo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &todo))
	assert.True(t, todo.Completed)

	rec = doRequest(e, http.MethodPatch, "/api/todos/"+itoa(created.ID)+"/toggle", nil)
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &todo))
	assert.False(t, todo.Completed)
}

func TestTodoHandler_Toggle_NotFound(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)

	rec := doRequest(e, http.MethodPatch, "/api/todos/999/toggle", nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTodoHandler_Delete(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	created := createTodo(t, e, "Delete me")

	rec := doRequest(e, http.MethodDelete, "/api/todos/"+itoa(created.ID), nil)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	rec = doRequest(e, http.MethodGet, "/api/todos/"+itoa(created.ID), nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTodoHandler_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)

	rec := doRequest(e, http.MethodDelete, "/api/todos/999", nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

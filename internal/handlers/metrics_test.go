package handlers_test

import (
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todoapp/internal/handlers"
	"todoapp/internal/metrics"
	"todoapp/internal/models"
)

// The metrics package registers its collectors once, process-wide, so these
// tests assert on deltas rather than absolute values — that way they're
// unaffected by whatever other tests ran before them in the same binary.

func TestMetrics_Create(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)

	createdBefore := testutil.ToFloat64(metrics.TodosCreatedTotal)
	openBefore := testutil.ToFloat64(metrics.TodosOpen)

	rec := doRequest(e, http.MethodPost, "/api/todos", handlers.CreateTodoRequest{Title: "Track me"})
	require.Equal(t, http.StatusCreated, rec.Code)

	assert.Equal(t, createdBefore+1, testutil.ToFloat64(metrics.TodosCreatedTotal))
	assert.Equal(t, openBefore+1, testutil.ToFloat64(metrics.TodosOpen))
}

func TestMetrics_Toggle_MovesOpenToClosedAndBack(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	created := createTodo(t, e, "Toggle me")
	todoPath := "/api/todos/" + itoa(created.ID) + "/toggle"

	openBefore := testutil.ToFloat64(metrics.TodosOpen)
	closedBefore := testutil.ToFloat64(metrics.TodosClosed)

	rec := doRequest(e, http.MethodPatch, todoPath, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, openBefore-1, testutil.ToFloat64(metrics.TodosOpen))
	assert.Equal(t, closedBefore+1, testutil.ToFloat64(metrics.TodosClosed))

	rec = doRequest(e, http.MethodPatch, todoPath, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, openBefore, testutil.ToFloat64(metrics.TodosOpen))
	assert.Equal(t, closedBefore, testutil.ToFloat64(metrics.TodosClosed))
}

func TestMetrics_Update_NoOpDoesNotMoveGauges(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	created := createTodo(t, e, "Leave me open")

	openBefore := testutil.ToFloat64(metrics.TodosOpen)
	closedBefore := testutil.ToFloat64(metrics.TodosClosed)

	rec := doRequest(e, http.MethodPut, "/api/todos/"+itoa(created.ID), handlers.UpdateTodoRequest{
		Title:     "Still open",
		Completed: false,
	})
	require.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, openBefore, testutil.ToFloat64(metrics.TodosOpen))
	assert.Equal(t, closedBefore, testutil.ToFloat64(metrics.TodosClosed))
}

func TestMetrics_Update_ClosingMovesOpenToClosed(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	created := createTodo(t, e, "Close me via update")

	openBefore := testutil.ToFloat64(metrics.TodosOpen)
	closedBefore := testutil.ToFloat64(metrics.TodosClosed)

	rec := doRequest(e, http.MethodPut, "/api/todos/"+itoa(created.ID), handlers.UpdateTodoRequest{
		Title:     "Closed via update",
		Completed: true,
	})
	require.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, openBefore-1, testutil.ToFloat64(metrics.TodosOpen))
	assert.Equal(t, closedBefore+1, testutil.ToFloat64(metrics.TodosClosed))
}

func TestMetrics_Delete_OpenTodoDecrementsOpen(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	created := createTodo(t, e, "Delete me while open")

	openBefore := testutil.ToFloat64(metrics.TodosOpen)
	closedBefore := testutil.ToFloat64(metrics.TodosClosed)

	rec := doRequest(e, http.MethodDelete, "/api/todos/"+itoa(created.ID), nil)
	require.Equal(t, http.StatusNoContent, rec.Code)

	assert.Equal(t, openBefore-1, testutil.ToFloat64(metrics.TodosOpen))
	assert.Equal(t, closedBefore, testutil.ToFloat64(metrics.TodosClosed))
}

func TestMetrics_Delete_ClosedTodoDecrementsClosed(t *testing.T) {
	db := setupTestDB(t)
	e := setupRouter(db)
	created := createTodo(t, e, "Delete me while closed")
	require.Equal(t, http.StatusOK, doRequest(e, http.MethodPatch, "/api/todos/"+itoa(created.ID)+"/toggle", nil).Code)

	openBefore := testutil.ToFloat64(metrics.TodosOpen)
	closedBefore := testutil.ToFloat64(metrics.TodosClosed)

	rec := doRequest(e, http.MethodDelete, "/api/todos/"+itoa(created.ID), nil)
	require.Equal(t, http.StatusNoContent, rec.Code)

	assert.Equal(t, openBefore, testutil.ToFloat64(metrics.TodosOpen))
	assert.Equal(t, closedBefore-1, testutil.ToFloat64(metrics.TodosClosed))
}

func TestMetrics_SeedFromDB(t *testing.T) {
	db := setupTestDB(t)
	require.NoError(t, db.Create(&models.Todo{Title: "Open 1", Completed: false}).Error)
	require.NoError(t, db.Create(&models.Todo{Title: "Open 2", Completed: false}).Error)
	require.NoError(t, db.Create(&models.Todo{Title: "Closed 1", Completed: true}).Error)

	require.NoError(t, metrics.SeedFromDB(db))

	assert.Equal(t, float64(2), testutil.ToFloat64(metrics.TodosOpen))
	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.TodosClosed))
}

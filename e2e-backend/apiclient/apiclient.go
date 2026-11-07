// Package apiclient is a thin HTTP client for the Todo app's backend, used
// by the E2E specs to exercise the real API over the network — as opposed to
// the backend's own unit tests, which call handlers in-process.
package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Todo struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ErrorBody struct {
	Error string `json:"error"`
}

// Response wraps a raw HTTP response so specs can assert on the status code
// and, separately, decode the body into whatever shape they expect.
type Response struct {
	StatusCode int
	Body       []byte
}

func (r *Response) DecodeTodo() (Todo, error) {
	var todo Todo
	err := json.Unmarshal(r.Body, &todo)
	return todo, err
}

func (r *Response) DecodeTodos() ([]Todo, error) {
	var todos []Todo
	err := json.Unmarshal(r.Body, &todos)
	return todos, err
}

func (r *Response) DecodeError() (ErrorBody, error) {
	var body ErrorBody
	err := json.Unmarshal(r.Body, &body)
	return body, err
}

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Healthz(ctx context.Context) (*Response, error) {
	return c.do(ctx, http.MethodGet, "/healthz", nil)
}

func (c *Client) Readyz(ctx context.Context) (*Response, error) {
	return c.do(ctx, http.MethodGet, "/readyz", nil)
}

func (c *Client) ListTodos(ctx context.Context) (*Response, error) {
	return c.do(ctx, http.MethodGet, "/api/todos", nil)
}

func (c *Client) CreateTodo(ctx context.Context, title, description string) (*Response, error) {
	return c.do(ctx, http.MethodPost, "/api/todos", map[string]string{
		"title":       title,
		"description": description,
	})
}

func (c *Client) GetTodo(ctx context.Context, id uint) (*Response, error) {
	return c.do(ctx, http.MethodGet, fmt.Sprintf("/api/todos/%d", id), nil)
}

func (c *Client) UpdateTodo(ctx context.Context, id uint, title, description string, completed bool) (*Response, error) {
	return c.do(ctx, http.MethodPut, fmt.Sprintf("/api/todos/%d", id), map[string]any{
		"title":       title,
		"description": description,
		"completed":   completed,
	})
}

func (c *Client) ToggleTodo(ctx context.Context, id uint) (*Response, error) {
	return c.do(ctx, http.MethodPatch, fmt.Sprintf("/api/todos/%d/toggle", id), nil)
}

func (c *Client) DeleteTodo(ctx context.Context, id uint) (*Response, error) {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/todos/%d", id), nil)
}

func (c *Client) do(ctx context.Context, method, path string, body any) (*Response, error) {
	var reqBody io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
		reqBody = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return &Response{StatusCode: resp.StatusCode, Body: data}, nil
}

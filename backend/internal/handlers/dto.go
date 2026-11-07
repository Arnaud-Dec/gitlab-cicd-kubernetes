package handlers

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type errorResponse struct {
	Error string `json:"error"`
}

package model

type TodoList struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoList struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

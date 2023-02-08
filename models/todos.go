package models

type Todo struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	IsComplete bool   `json:"complete"`
}

type TodoList struct {
	Todos []Todo `json:"TodoList"`
}

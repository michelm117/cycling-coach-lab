package models

type Task struct {
	ID int `json:"id"`
	Title  string `json:"title"`
	Content string `json:"content"`
	// Completed bool `json:"completed"`
}

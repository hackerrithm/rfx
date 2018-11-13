package domain

import "time"

type (
	// Task is the struct that would contain any
	// domain logic if we had any. Because it's simple
	// we're going to send it over the wire directly
	// so we add the JSON serialization tags but we
	// could use DTO specific structs for that
	Task struct {
		ID      int64     `json:"id" bson:"_id"`
		Author  string    `json:"author" bson:"author,omitempty"`
		Content string    `json:"content" bson:"content,omitempty"`
		Date    time.Time `json:"timestamp" bson:"timestamp,omitempty"`
	}
)

// NewTask creates a new Task
func NewTask(author, content string) *Task {
	return &Task{
		Author:  author,
		Content: content,
		Date:    now(),
	}
}

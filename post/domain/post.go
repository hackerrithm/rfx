package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

var postContextKey contextKey = "post"

type (
	contextKey string

	// Post is the struct that would contain any
	// domain logic if we had any. Because it's simple
	// we're going to send it over the wire directly
	// so we add the JSON serialization tags but we
	// could use DTO specific structs for that
	Post struct {
		ID           bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Author       string        `json:"author" bson:"author"`
		Topic        string        `json:"topic" bson:"topic"`
		Category     string        `json:"category" bson:"category"`
		ContentText  string        `json:"contentText" bson:"content_text"`
		ContentPhoto string        `json:"contentPhoto" bson:"content_photo"`
		Date         time.Time     `json:"timestamp" bson:"timestamp"`
	}
)

// NewPost creates a new Post!
func NewPost(author, topic, category, text, photo string) *Post {
	return &Post{
		Author:       author,
		Topic:        topic,
		Category:     category,
		ContentText:  text,
		ContentPhoto: photo,
		Date:         now(),
	}
}

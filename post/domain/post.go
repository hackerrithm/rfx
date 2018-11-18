package domain

import "time"

var postContextKey contextKey = "Post"

type (
	contextKey string

	// Post is the struct that would contain any
	// domain logic if we had any. Because it's simple
	// we're going to send it over the wire directly
	// so we add the JSON serialization tags but we
	// could use DTO specific structs for that
	Post struct {
		ID           string    `json:"_id"`
		Author       string    `json:"author"`
		Topic        string    `json:"topic"`
		Category     string    `json:"category"`
		ContentText  string    `json:"contentText"`
		ContentPhoto string    `json:"contentPhoto"`
		Date         time.Time `json:"timestamp"`
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

// // SetPassword sets Post's password
// func (u *Post) SetPassword(p string) {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
// 	if err != nil {
// 		panic(err)
// 	}

// 	u.Password = string(hashedPassword)
// }

// // IsCredentialsVerified ...
// func (u *Post) IsCredentialsVerified(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

// // NewContext creates a new context
// func (u *Post) NewContext(ctx context.Context) context.Context {
// 	return context.WithValue(ctx, postContextKey, u)
// }

// // PostFromContext gets Post from context
// func PostFromContext(ctx context.Context) (*Post, bool) {
// 	u, ok := ctx.Value(postContextKey).(*Post)
// 	return u, ok
// }

// // PostMustFromContext gets Post from context. if can't make panic
// func PostMustFromContext(ctx context.Context) *Post {
// 	u, ok := ctx.Value(postContextKey).(*Post)
// 	if !ok {
// 		panic("Post can't get from request's context")
// 	}
// 	return u
// }

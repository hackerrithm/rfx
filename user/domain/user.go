package domain

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var userContextKey contextKey = "user"

type (
	contextKey string

	// User is the struct that would contain any
	// domain logic if we had any. Because it's simple
	// we're going to send it over the wire directly
	// so we add the JSON serialization tags but we
	// could use DTO specific structs for that
	User struct {
		ID        string    `json:"_id"`
		UserName  string    `json:"username"`
		Password  string    `json:"paswword"`
		FirstName string    `json:"firstname"`
		LastName  string    `json:"lastname"`
		Gender    string    `json:"gender"`
		Date      time.Time `json:"timestamp"`
	}
)

// NewUser creates a new User!
func NewUser(userName, password, firstname, lastname, gender string) *User {
	return &User{
		UserName:  userName,
		Password:  password,
		FirstName: firstname,
		LastName:  lastname,
		Gender:    gender,
		Date:      now(),
	}
}

// SetPassword sets user's password
func (u *User) SetPassword(p string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	u.Password = string(hashedPassword)
}

// IsCredentialsVerified ...
func (u *User) IsCredentialsVerified(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// NewContext creates a new context
func (u *User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// UserFromContext gets user from context
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	return u, ok
}

// UserMustFromContext gets user from context. if can't make panic
func UserMustFromContext(ctx context.Context) *User {
	u, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		panic("user can't get from request's context")
	}
	return u
}

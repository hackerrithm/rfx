package listing

import "errors"

// ErrNotFound is used when a user could not be found.
var ErrNotFound = errors.New("user not found")

// Repository provides access to the users.
type Repository interface {
	// GetUser returns the user with given ID.
	GetUser(int) (User, error)
	// GetAllUsers returns all beers saved in storage.
	GetAllUsers() []User
}

// Service provides user and review listing operations.
type Service interface {
	GetUser(int) (User, error)
	GetUsers() []User
}

type service struct {
	r Repository
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// GetUsers returns all users
func (s *service) GetUsers() []User {
	return s.r.GetAllUsers()
}

// GetUser returns a user
func (s *service) GetUser(id int) (User, error) {
	return s.r.GetUser(id)
}

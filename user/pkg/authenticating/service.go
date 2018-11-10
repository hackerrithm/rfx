package authenticating

import "errors"

// Service ...
type Service interface {
	Login(string, string) (interface{}, error)
	SignUp(string, string, string, string) (interface{}, error)
}

// ErrDuplicate ...
var ErrDuplicate = errors.New("user already exists")

// Repository provides access to user repository.
type Repository interface {
	// Login returns a token if user is in storage.
	Login(string, string) (interface{}, error)
	// SignUp inserts a new user and returns a token if user is in storage.
	SignUp(string, string, string, string) (interface{}, error)
}

type service struct {
	uSR Repository
}

// NewService ...
func NewService(r Repository) Service {
	return &service{r}
}

// Login ...
func (s *service) Login(username string, password string) (interface{}, error) {
	return s.uSR.Login(username, password)
}

// SignUp ...
func (s *service) SignUp(username string, password string, firstname string, lastname string) (interface{}, error) {
	return s.uSR.SignUp(username, password, firstname, lastname)
}

package authenticating

import "errors"

// ErrDuplicate ...
var ErrDuplicate = errors.New("user already exists")

// Service ...
type (
	Service interface {
		Login(string, string) (map[string]interface{}, error)
		SignUp(string, string, string, string) (interface{}, error)
		Profile(string, string) ([]byte, error)
	}

	// Repository provides access to user repository.
	Repository interface {
		// Login returns a token if user is in storage.
		Login(string, string) (map[string]interface{}, error)
		// SignUp inserts a new user and returns a token if user is in storage.
		SignUp(string, string, string, string) (interface{}, error)
		// Profile ...
		Profile(string, string) ([]byte, error)
	}

	service struct {
		usr Repository
	}
)

// NewService ...
func NewService(r Repository) Service {
	return &service{r}
}

// Login ...
func (s *service) Login(username string, password string) (map[string]interface{}, error) {
	return s.usr.Login(username, password)
}

// SignUp ...
func (s *service) SignUp(username string, password string, firstname string, lastname string) (interface{}, error) {
	return s.usr.SignUp(username, password, firstname, lastname)
}

// SignUp ...
func (s *service) Profile(token string, UUID string) ([]byte, error) {
	return s.usr.Profile(token, UUID)
}

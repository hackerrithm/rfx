package authenticating

import "errors"

type (
	// Payload ...
	Payload []User

	// Event ...
	Event int
)

const (
	// Done means finished processing successfully
	Done Event = iota

	// UserAlreadyExists means the given user is a duplicate of an existing one
	UserAlreadyExists

	// Failed means processing did not finish successfully
	Failed
)

// GetMeaning ...
func (e Event) GetMeaning() string {
	if e == Done {
		return "Done"
	}

	if e == UserAlreadyExists {
		return "User Already exists"
	}

	if e == Failed {
		return "Failed"
	}

	return "Unknown result"
}

// ErrDuplicate ...
var ErrDuplicate = errors.New("user already exists")

// Service ...
type Service interface {
	AddUser(...User)
	AddSampleUsers(Payload) <-chan Event
	Login(string, string) ([]byte, error)
}

// Repository provides access to user repository.
type Repository interface {
	// AddUser saves a given user to the repository.
	AddUser(User) error
	// Login returns a token if user is in storage.
	Login(string, string) ([]byte, error)
}

type service struct {
	uSR Repository
}

// NewService ...
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddUser(u ...User) {
	for _, user := range u {
		_ = s.uSR.AddUser(user)
	}
}

// AddSampleUsers adds some sample users to the database
func (s *service) AddSampleUsers(data Payload) <-chan Event {
	results := make(chan Event)

	go func() {
		defer close(results)

		for _, b := range data {
			err := s.uSR.AddUser(b)
			if err != nil {
				if err == ErrDuplicate {
					// forgive the naughty error type checking above...
					results <- UserAlreadyExists
					continue
				}
				results <- Failed
				continue
			}

			results <- Done
		}
	}()

	return results
}

// Login ...
func (s *service) Login(username string, password string) ([]byte, error) {
	return s.uSR.Login(username, password)
}

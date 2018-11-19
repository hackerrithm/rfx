package engine

import (
	"sync"

	"golang.org/x/net/context"
)

type (
	// User is the interface for our interactor
	User interface {
		// Add is the add-a-user use-case
		Add(c context.Context, r *AddUserRequest) *AddUserResponse

		// Update is the update-a-user use-case
		Update(c context.Context, r *AddUserRequest) *UpdateUserResponse

		// List is the list-the-users use-case
		List(c context.Context, r *ListUsersRequest) *ListUsersResponse

		// Read is the authenticate use-case
		Read(c context.Context, r *AuthenticateUserRequest) *AuthenticatUserResponse

		// Profile is the getting profile details use-case
		Profile(c context.Context, r *ProfileRequest) *ProfileResponse

		//GenerateToken ...
		GenerateToken() (map[string]interface{}, error)

		// ParseToken ...
		ParseToken(token string) (map[string]interface{}, error)
	}

	// user implementation
	user struct {
		repository UserRepository
		jwt        JWTSignParser
	}
)

var (
	userInstance User
	userOnce     sync.Once
)

// NewUser creates a new User interactor wired up
// to use the user repository from the storage provider
// that the engine has been setup to use.
func (f *engineFactory) NewUser() User {
	userOnce.Do(func() {
		userInstance = &user{
			repository: f.NewUserRepository(),
			jwt:        f.JWTSignParser,
		}
	})
	return userInstance
}

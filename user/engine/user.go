package engine

import (
	"golang.org/x/net/context"
)

type (
	// User is the interface for our interactor
	User interface {
		// Add is the add-a-user use-case
		Add(c context.Context, r *AddUserRequest) *AddUserResponse

		// List is the list-the-users use-case
		List(c context.Context, r *ListUsersRequest) *ListUsersResponse

		// Read is the authenticate use-case
		Read(c context.Context, r *AuthenticateUserRequest) *AuthenticatUserResponse

		// Profile is the getting profile details use-case
		Profile(c context.Context, r *ProfileRequest) *ProfileResponse
	}

	// user implementation
	user struct {
		repository UserRepository
	}
)

// NewUser creates a new User interactor wired up
// to use the user repository from the storage provider
// that the engine has been setup to use.
func (f *engineFactory) NewUser() User {
	return &user{
		repository: f.NewUserRepository(),
	}
}

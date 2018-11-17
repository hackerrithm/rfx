package engine

import (
	"golang.org/x/net/context"

	"github.com/hackerrithm/longterm/rfx/user/domain"
)

type (
	// UserRepository defines the methods that any
	// data storage provider needs to implement to get
	// and store users
	UserRepository interface {
		// Insert adds a new User to the datastore
		Insert(c context.Context, user *domain.User) string

		// Put adds a new User to the datastore
		Put(c context.Context, user *domain.User) string

		// List returns existing users matching the
		// query provided
		List(c context.Context, query *Query) []*domain.User

		// Read returns ...
		Read(c context.Context, username, password string) *domain.User

		// Profile returns user details on a specific profile
		Profile(c context.Context, jwtToken, ID string) *domain.User
	}

	// StorageFactory is the interface that a storage
	// provider needs to implement so that the engine can
	// request repository instances as it needs them
	StorageFactory interface {
		// NewUserRepository returns a storage specific
		// UserRepository implementation
		NewUserRepository() UserRepository
	}
)

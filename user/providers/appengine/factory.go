// +build appengine

package appengine

import (
	"github.com/hackerrithm/longterm/rfx/user/engine"
)

type (
	storageFactory struct{}
)

// NewStorage creates a new instance of this datastore storage factory
func NewStorage() engine.StorageFactory {
	return &storageFactory{}
}

// NewUserRepository creates a new datastore greeting repository
func (f *storageFactory) NewUserRepository() engine.UserRepository {
	return newUserRepository()
}

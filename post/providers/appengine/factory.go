// +build appengine

package appengine

import (
	"github.com/hackerrithm/longterm/rfx/post/engine"
)

type (
	storageFactory struct{}
)

// NewStorage creates a new instance of this datastore storage factory
func NewStorage() engine.StorageFactory {
	return &storageFactory{}
}

// NewPostRepository creates a new datastore greeting repository
func (f *storageFactory) NewPostRepository() engine.PostRepository {
	return newPostRepository()
}

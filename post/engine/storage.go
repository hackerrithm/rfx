package engine

import (
	"golang.org/x/net/context"

	"github.com/hackerrithm/longterm/rfx/post/domain"
)

type (
	// PostRepository defines the methods that any
	// data storage provider needs to implement to get
	// and store posts
	PostRepository interface {
		// Insert adds a new Post to the datastore
		Insert(c context.Context, post *domain.Post) string

		// Put adds a new Post to the datastore
		Put(c context.Context, post *domain.Post) string

		// List returns existing posts matching the
		// query provided
		List(c context.Context, query *Query) []*domain.Post

		// Read returns ...
		Read(c context.Context, id string) *domain.Post
	}

	// StorageFactory is the interface that a storage
	// provider needs to implement so that the engine can
	// request repository instances as it needs them
	StorageFactory interface {
		// NewPostRepository returns a storage specific
		// PostRepository implementation
		NewPostRepository() PostRepository
	}
)

package engine

import (
	"context"

	"github.com/hackerrithm/longterm/rfx/task/domain"
)

type (
	// TaskRepository defines the methods that any
	// data storage provider needs to implement to get
	// and store Tasks
	TaskRepository interface {
		// Put adds a new Task to the datastore
		Put(c context.Context, Task *domain.Task)

		// List returns existing Tasks matching the
		// query provided
		List(c context.Context, query *Query) []*domain.Task
	}

	// StorageFactory is the interface that a storage
	// provider needs to implement so that the engine can
	// request repository instances as it needs them
	StorageFactory interface {
		// NewTaskRepository returns a storage specific
		// TaskRepository implementation
		NewTaskRepository() TaskRepository
	}
)

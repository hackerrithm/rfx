package engine

import "context"

type (
	// Task is the interface for our interactor
	Task interface {
		// Add is the add-a-Task use-case
		Add(c context.Context, r *AddTaskRequest) *AddTaskResponse

		// List is the list-the-Tasks use-case
		List(c context.Context, r *ListTasksRequest) *ListTasksResponse
	}

	// task implementation
	task struct {
		repository TaskRepository
	}
)

// NewTask creates a new Task interactor wired up
// to use the task repository from the storage provider
// that the engine has been setup to use.
func (f *engineFactory) NewTask() Task {
	return &task{
		repository: f.NewTaskRepository(),
	}
}

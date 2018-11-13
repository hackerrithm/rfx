package engine

import (
	"context"

	"github.com/hackerrithm/longterm/rfx/task/domain"
)

type (
	// AddTaskRequest ...
	AddTaskRequest struct {
		Author  string
		Content string
	}

	// AddTaskResponse ...
	AddTaskResponse struct {
		ID int64
	}
)

func (t *task) Add(c context.Context, r *AddTaskRequest) *AddTaskResponse {
	// this is where all our app logic would go - the
	// rules that apply to adding a Task whether it
	// is being done via the web UI, a console app, or
	// whatever the internet has just been added to ...
	Task := domain.NewTask(r.Author, r.Content)
	t.repository.Put(c, Task)
	return &AddTaskResponse{
		ID: Task.ID,
	}
}

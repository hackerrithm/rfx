package engine

import (
	"context"

	"github.com/hackerrithm/longterm/rfx/task/domain"
)

type (
	// ListTasksRequest ...
	ListTasksRequest struct {
		Count int
	}

	// ListTasksResponse ...
	ListTasksResponse struct {
		Tasks []*domain.Task
	}
)

func (t *task) List(c context.Context, r *ListTasksRequest) *ListTasksResponse {
	q := NewQuery("Task").Order("date", Descending).Slice(0, r.Count)
	return &ListTasksResponse{
		Tasks: t.repository.List(c, q),
	}
}

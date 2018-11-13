package engine

import (
	"context"

	"github.com/hackerrithm/longterm/rfx/user/domain"
)

type (
	// ListUsersRequest ...
	ListUsersRequest struct {
		Count int
	}

	// ListUsersResponse ...
	ListUsersResponse struct {
		Users []*domain.User
	}
)

func (u *user) List(c context.Context, r *ListUsersRequest) *ListUsersResponse {
	q := NewQuery("User").Order("date", Descending).Slice(0, r.Count)
	return &ListUsersResponse{
		Users: u.repository.List(c, q),
	}
}

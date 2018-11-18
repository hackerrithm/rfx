package engine

import (
	"context"

	"github.com/hackerrithm/longterm/rfx/post/domain"
)

type (
	// ListPostsRequest ...
	ListPostsRequest struct {
		Count int
	}

	// ListPostsResponse ...
	ListPostsResponse struct {
		Posts []*domain.Post
	}
)

func (p *post) List(c context.Context, r *ListPostsRequest) *ListPostsResponse {
	q := NewQuery("Post").Order("date", Descending).Slice(0, r.Count)
	return &ListPostsResponse{
		Posts: p.repository.List(c, q),
	}
}

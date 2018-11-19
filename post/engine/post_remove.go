package engine

import (
	"golang.org/x/net/context"
)

type (
	// RemovePostRequest ...
	RemovePostRequest struct {
		ID string
	}

	// RemovePostResponse ...
	RemovePostResponse struct {
		Post string
	}
)

func (p *post) Remove(c context.Context, r *RemovePostRequest) *RemovePostResponse {
	return &RemovePostResponse{
		p.repository.Remove(c, r.ID),
	}
}

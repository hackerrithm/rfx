package engine

import (
	"golang.org/x/net/context"

	"github.com/hackerrithm/longterm/rfx/post/domain"
)

type (
	// AddPostRequest ...
	AddPostRequest struct {
		Author       string
		Topic        string
		Category     string
		ContentText  string
		ContentPhoto string
	}

	// AddPostResponse ...
	AddPostResponse struct {
		ID string
	}
)

func (p *post) Add(c context.Context, r *AddPostRequest) *AddPostResponse {
	Post := domain.NewPost(r.Author, r.Topic, r.Category, r.ContentText, r.ContentPhoto)
	return &AddPostResponse{
		ID: p.repository.Insert(c, Post),
	}
}

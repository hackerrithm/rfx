package engine

import (
	"golang.org/x/net/context"

	"github.com/hackerrithm/longterm/rfx/post/domain"
)

type (
	// UpdatePostRequest ...
	UpdatePostRequest struct {
		Author       string
		Topic        string
		Category     string
		ContentText  string
		ContentPhoto string
	}

	// UpdatePostResponse ...
	UpdatePostResponse struct {
		ID string
	}
)

func (p *post) Update(c context.Context, r *UpdatePostRequest) *UpdatePostResponse {
	// this is where all our app logic would go - the
	// rules that apply to adding a Post whether it
	// is being done via the web UI, a console app, or
	// whatever the internet has just been added to ...
	post := domain.NewPost(r.Author, r.Topic, r.Category, r.ContentText, r.ContentPhoto)
	// u.repository.Put(c, Post)
	return &UpdatePostResponse{
		ID: p.repository.Put(c, post),
	}
}

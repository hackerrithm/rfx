package engine

import (
	"golang.org/x/net/context"
)

type (
	// Post is the interface for our interactor
	Post interface {
		// Add is the add-a-Post use-case
		Add(c context.Context, r *AddPostRequest) *AddPostResponse

		// Update is the update-a-post use-case
		Update(c context.Context, r *UpdatePostRequest, id string) *UpdatePostResponse

		// List is the list-the-posts use-case
		List(c context.Context, r *ListPostsRequest) *ListPostsResponse

		// Read is the authenticate use-case
		Read(c context.Context, r *ReadPostRequest) *ReadPostResponse

		// Remove ...
		Remove(c context.Context, r *RemovePostRequest) *RemovePostResponse
	}

	// post implementation
	post struct {
		repository PostRepository
		// jwt        JWTSignParser
	}
)

// NewPost creates a new Post interactor wired up
// to use the Post repository from the storage provider
// that the engine has been setup to use.
func (f *engineFactory) NewPost() Post {
	return &post{
		repository: f.NewPostRepository(),
		// jwt:        f.JWTSignParser,
	}
}

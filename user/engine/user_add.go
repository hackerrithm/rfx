package engine

import (
	"golang.org/x/net/context"

	"github.com/hackerrithm/longterm/rfx/user/domain"
)

type (
	// AddUserRequest ...
	AddUserRequest struct {
		UserName  string
		Password  string
		FirstName string
		LastName  string
		Gender    string
	}

	// AddUserResponse ...
	AddUserResponse struct {
		ID string
	}
)

func (u *user) Add(c context.Context, r *AddUserRequest) *AddUserResponse {
	user := domain.NewUser(r.UserName, r.Password, r.FirstName, r.LastName, r.Gender)
	return &AddUserResponse{
		ID: u.repository.Insert(c, user),
	}
}

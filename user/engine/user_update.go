package engine

import (
	"golang.org/x/net/context"

	"github.com/hackerrithm/longterm/rfx/user/domain"
)

type (
	// UpdateUserRequest ...
	UpdateUserRequest struct {
		UserName  string
		Password  string
		FirstName string
		LastName  string
		Gender    string
	}

	// UpdateUserResponse ...
	UpdateUserResponse struct {
		Token string
	}
)

func (u *user) Update(c context.Context, r *AddUserRequest) *UpdateUserResponse {
	// this is where all our app logic would go - the
	// rules that apply to adding a user whether it
	// is being done via the web UI, a console app, or
	// whatever the internet has just been added to ...
	user := domain.NewUser(r.UserName, r.Password, r.FirstName, r.LastName, r.Gender)
	// u.repository.Put(c, user)
	return &UpdateUserResponse{
		Token: u.repository.Put(c, user),
	}
}

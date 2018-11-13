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
		Token map[string]interface{}
	}
)

func (u *user) Add(c context.Context, r *AddUserRequest) *AddUserResponse {
	// this is where all our app logic would go - the
	// rules that apply to adding a user whether it
	// is being done via the web UI, a console app, or
	// whatever the internet has just been added to ...
	user := domain.NewUser(r.UserName, r.Password, r.FirstName, r.LastName, r.Gender)
	// u.repository.Put(c, user)
	return &AddUserResponse{
		Token: u.repository.Put(c, user),
	}
}

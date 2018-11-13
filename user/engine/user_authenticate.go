package engine

import (
	"golang.org/x/net/context"
)

type (
	// AuthenticateUserRequest ...
	AuthenticateUserRequest struct {
		UserName string
		Password string
	}

	// AuthenticatUserResponse ...
	AuthenticatUserResponse struct {
		Token map[string]interface{}
	}

	// ProfileRequest ...
	ProfileRequest struct {
		Token string
		ID    string
	}

	// ProfileResponse ...
	ProfileResponse struct {
		Payload []byte
	}
)

func (u *user) Read(c context.Context, r *AuthenticateUserRequest) *AuthenticatUserResponse {
	// this is where all our app logic would go - the
	// rules that apply to adding a user whether it
	// is being done via the web UI, a console app, or
	// whatever the internet has just been added to ...
	// user := domain.NewUser(r.UserName, r.Password, r.FirstName, r.LastName, r.Gender)

	return &AuthenticatUserResponse{
		u.repository.Read(c, r.UserName, r.Password),
	}
}

func (u *user) Profile(c context.Context, r *ProfileRequest) *ProfileResponse {

	return &ProfileResponse{
		u.repository.Profile(c, r.Token, r.ID),
	}
}

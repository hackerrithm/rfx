package engine

import (
	"github.com/hackerrithm/longterm/rfx/user/domain"
	"golang.org/x/net/context"
)

const (
	secretKey = "12This98Is34A76String56Used65As78Secret01"
)

type (
	// AuthenticateUserRequest ...
	AuthenticateUserRequest struct {
		UserName string
		Password string
	}

	// AuthenticatUserResponse ...
	AuthenticatUserResponse struct {
		User *domain.User
	}

	// ProfileRequest ...
	ProfileRequest struct {
		Token string
		ID    string
	}

	// ProfileResponse ...
	ProfileResponse struct {
		Payload *domain.User
	}
)

func (u *user) Read(c context.Context, r *AuthenticateUserRequest) *AuthenticatUserResponse {
	return &AuthenticatUserResponse{
		u.repository.Read(c, r.UserName, r.Password),
	}
}

func (u *user) Profile(c context.Context, r *ProfileRequest) *ProfileResponse {
	return &ProfileResponse{
		u.repository.Profile(c, r.Token, r.ID),
	}
}

func (u *user) GenerateToken() (map[string]interface{}, error) {
	claims := map[string]interface{}{
		"userid": "u1",
	}

	return u.jwt.Sign(claims, secretKey)
}

func (u *user) ParseToken(token string) (map[string]interface{}, error) {
	return u.jwt.Parse(token, secretKey)
}

package engine

import (
	"github.com/hackerrithm/longterm/rfx/post/domain"
	"golang.org/x/net/context"
)

const (
	secretKey = "12This98Is34A76String56Used65As78Secret01"
)

type (
	// ReadPostRequest ...
	ReadPostRequest struct {
		ID string
	}

	// ReadPostResponse ...
	ReadPostResponse struct {
		Post *domain.Post
	}
)

func (p *post) Read(c context.Context, r *ReadPostRequest) *ReadPostResponse {
	return &ReadPostResponse{
		p.repository.Read(c, r.ID),
	}
}

// func (u *Post) GenerateToken() (map[string]interface{}, error) {
// 	claims := map[string]interface{}{
// 		"postid": "u1",
// 	}

// 	return u.jwt.Sign(claims, secretKey)
// }

// func (u *Post) ParseToken(token string) (map[string]interface{}, error) {
// 	return u.jwt.Parse(token, secretKey)
// }

package security

import (
	"log"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/hackerrithm/longterm/rfx/user/engine"
)

// NewJWT ...
func NewJWT() engine.JWTSignParser {
	return &jwt{}
}

type jwt struct{}

var returnObjectMap = make(map[string]interface{})

// func getByteToken(token *jwtgo.Token) (interface{}, error) {
// 	if jwtgo.SigningMethodHS256 != token.Method {
// 		log.Println("Invalid signing algorithm")
// 	}

// 	return []byte(secretKey), nil
// }

// func parseWithClaims(jwtToken string) (*jwtgo.Token, error) {
// 	cl, err := jwtgo.ParseWithClaims(jwtToken, &JWTData{}, getByteToken)
// 	if err != nil {
// 		log.Println("error in parseWithClaims")
// 	}
// 	return cl, nil

// }

// Sign ...
func (j *jwt) Sign(claims map[string]interface{}, secretKey string) (map[string]interface{}, error) {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, jwtgo.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("StatusUnauthorized ", err)
	}

	returnObjectMap["token"] = tokenString

	return returnObjectMap, nil
}

// Parse ..
func (j *jwt) Parse(tokenStr, secret string) (map[string]interface{}, error) {
	token, err := jwtgo.Parse(tokenStr, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		e, ok := err.(*jwtgo.ValidationError)
		if ok {
			return nil, e
		}
		return nil, err
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}

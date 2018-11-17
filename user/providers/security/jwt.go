package security

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// // NewJWT ...
// func NewJWT() engine.JWTSignParser {
// 	return &jwt{}
// }

// type jwt struct{}

var returnObjectMap map[string]interface{}

// JWTData is a struct with the structure of the jwt data
type JWTData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

const (
	secretKey = "12This98Is34A76String56Used65As78Secret01"
)

func getByteToken(token *jwt.Token) (interface{}, error) {
	if jwt.SigningMethodHS256 != token.Method {
		log.Println("Invalid signing algorithm")
	}

	return []byte(secretKey), nil
}

func parseWithClaims(jwtToken string) (*jwt.Token, error) {
	cl, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, getByteToken)
	if err != nil {
		log.Println("error in parseWithClaims")
	}
	return cl, nil

}

// Sign ...
func Sign( /*claims map[string]interface{}, secret string*/ ) (map[string]interface{}, error) {
	returnObjectMap = make(map[string]interface{})
	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},

		CustomClaims: map[string]string{
			"userid": "u1",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("StatusUnauthorized ", err)
	}

	returnObjectMap["token"] = tokenString

	return returnObjectMap, nil
}

// Parse ..
func Parse(tokenStr string /*secret string, , userUID int64*/) (string, error) {
	claims, err := parseWithClaims(tokenStr)
	if err != nil {
		log.Println(err)
		return "", err
	}

	data := claims.Claims.(*JWTData)

	userID := data.CustomClaims["userid"]
	log.Println("claim ", userID)

	return userID, nil
}

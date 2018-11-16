package security

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	cfg "github.com/hackerrithm/longterm/rfx/configs"
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

func getSecret() (string, error) {
	file, _ := os.Open("../configs/config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := cfg.Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	return configuration.SECRET, nil
}

func getByteToken(token *jwt.Token) (interface{}, error) {
	secretKey, err := getSecret()
	if err != nil {
		return nil, err
	}
	if jwt.SigningMethodHS256 != token.Method {
		log.Println("Invalid signing algorithm")
	}

	return []byte(secretKey), nil
}

func ParseWithClaims(jwtToken string) (*jwt.Token, error) {
	cl, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, getByteToken)
	if err != nil {
		log.Println("error in parseWithClaims")
	}
	return cl, nil

}

func Sign( /*claims map[string]interface{}, secret string*/ ) (map[string]interface{}, error) {
	returnObjectMap = make(map[string]interface{})
	var result []byte
	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},

		CustomClaims: map[string]string{
			"userid": "u1",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey, err := getSecret()
	if err != nil {
		return nil, err
	}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("StatusUnauthorized ", err)
	}

	result, err = json.Marshal(struct {
		Token string `json:"token"`
	}{
		tokenString,
	})

	returnObjectMap["token"] = result

	return returnObjectMap, nil
}

func Parse( /*tokenStr string, secret string, */ userUID int64) (map[string]interface{}, error) {
	returnObjectMap = make(map[string]interface{})
	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},

		CustomClaims: map[string]string{
			"userid": string(userUID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey, err := getSecret()
	if err != nil {
		return nil, err
	}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("StatusUnauthorized ", err)
		return nil, err
	}

	returnObjectMap["token"] = tokenString
	returnObjectMap["userUID"] = userUID

	return returnObjectMap, nil
}

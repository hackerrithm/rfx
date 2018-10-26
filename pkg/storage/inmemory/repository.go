package inmemory

import (
	"encoding/json"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hackerrithm/longterm/rfx/internal/pkg/authenticating"
)

// Storage (Memory storage) keeps data in memory
type Storage struct {
	users []User
}

// AddUser saves the given user to the repository
func (m *Storage) AddUser(u authenticating.User) error {
	for _, e := range m.users {
		if u.UserName == e.UserName &&
			u.FirstName == e.FirstName &&
			u.LastName == e.LastName {
			return authenticating.ErrDuplicate
		}
	}

	newU := User{
		UID:       len(m.users) + 1,
		UserName:  u.UserName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Gender:    u.Gender,
		Password:  u.Password,
	}
	m.users = append(m.users, newU)

	return nil
}

// TODO: to be implemented

// JWTData is a struct with the structure of the jwt data
type JWTData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

// Login ...
func (m *Storage) Login(username string, password string) (interface{}, error) {

	var user authenticating.User
	var result []byte

	user.UserName = username
	user.Password = password

	// Demo - in real case scenario you'd check this against your database
	if user.UserName == "admin" && user.Password == "password" {
		claims := JWTData{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
			},

			CustomClaims: map[string]string{
				"userid": "u1",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte( /*SECRET*/ "asdad"))
		if err != nil {
			log.Println("StatusUnauthorized ", err)
		}

		result, err = json.Marshal(struct {
			Token string `json:"token"`
		}{
			tokenString,
		})

		if err != nil {
			log.Println("StatusUnauthorized ", err)
		}
	}
	return result, nil
}

// SignUp ...
func (m *Storage) SignUp(username string, password string, firstname string, lastname string) (interface{}, error) {

	var user authenticating.User
	var result []byte

	user.UserName = username
	user.Password = password
	user.FirstName = firstname
	user.LastName = lastname
	user.Gender = ""

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},

		CustomClaims: map[string]string{
			"userid": "u1",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte( /*SECRET*/ "asdad"))
	if err != nil {
		log.Println("StatusUnauthorized ", err)
	}

	result, err = json.Marshal(struct {
		Token string `json:"token"`
	}{
		tokenString,
	})

	if err != nil {
		log.Println("StatusUnauthorized ", err)
	}
	return result, nil

}

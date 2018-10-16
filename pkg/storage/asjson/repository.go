package asjson

import (
	"encoding/json"
	"log"
	"path"
	"runtime"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hackerrithm/longterm/rfx/internal/pkg/authenticating"
	"github.com/hackerrithm/longterm/rfx/internal/pkg/listing"
	scribble "github.com/nanobox-io/golang-scribble"
)

const (
	// dir defines the name of the directory where the files are stored
	dir = "/data/"

	// CollectionUser identifier for the JSON collection of users
	CollectionUser = "users"
)

// Storage stores user data in JSON files
type Storage struct {
	db *scribble.Driver
}

// NewStorage returns a new JSON  storage
func NewStorage() (*Storage, error) {
	var err error

	s := new(Storage)

	_, filename, _, _ := runtime.Caller(0)
	p := path.Dir(filename)

	s.db, err = scribble.New(p+dir, nil)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// AddUser saves the given user to the repository
func (s *Storage) AddUser(u authenticating.User) error {

	existingUsers := s.GetAllUsers()
	for _, e := range existingUsers {
		if u.UserName == e.UserName &&
			u.FirstName == e.FirstName &&
			u.LastName == e.LastName {
			return authenticating.ErrDuplicate
		}
	}

	newU := User{
		UID:       len(existingUsers) + 1,
		UserName:  u.UserName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Gender:    u.Gender,
		Password:  u.Password,
	}
	// m.users = append(m.users, newU)

	resource := strconv.Itoa(newU.UID)
	if err := s.db.Write(CollectionUser, resource, newU); err != nil {
		return err
	}
	return nil
}

// GetUser returns a user with the specified ID
func (s *Storage) GetUser(id int) (listing.User, error) {
	var b User
	var user listing.User

	var resource = strconv.Itoa(id)

	if err := s.db.Read(CollectionUser, resource, &b); err != nil {
		// err handling omitted for simplicity
		return user, listing.ErrNotFound
	}

	user.UID = b.UID
	user.UserName = b.UserName
	user.FirstName = b.FirstName
	user.LastName = b.LastName
	user.Password = b.Password
	user.Gender = b.Gender

	return user, nil
}

// GetAllUsers returns all users
func (s *Storage) GetAllUsers() []listing.User {
	list := []listing.User{}

	records, err := s.db.ReadAll(CollectionUser)
	if err != nil {
		// err handling omitted for simplicity
		return list
	}

	for _, r := range records {
		var b User
		var user listing.User

		if err := json.Unmarshal([]byte(r), &b); err != nil {
			// err handling omitted for simplicity
			return list
		}

		user.UID = b.UID
		user.UserName = b.UserName
		user.FirstName = b.FirstName
		user.LastName = b.LastName
		user.Password = b.Password
		user.Gender = b.Gender

		list = append(list, user)
	}

	return list
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
func (s *Storage) Login(username string, password string) ([]byte, error) {

	var user listing.User
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

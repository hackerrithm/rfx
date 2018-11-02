package mongo

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hackerrithm/longterm/rfx/internal/pkg/authenticating"
	"github.com/hackerrithm/longterm/rfx/internal/pkg/listing"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Storage ...
type Storage struct {
	Database string
}

var db *mgo.Database

var returnObjectMap map[string]interface{}

const (
	// COLLECTION ...
	COLLECTION = "user"
)

// NewStorage ...
func NewStorage() (*Storage, error) {
	session, err := mgo.Dial("mongodb://localhost/test1")
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB("test1")
	return &Storage{db.Name}, err
}

// GetAllUsers ...
func (s *Storage) GetAllUsers() []listing.User {
	var movies []listing.User
	err := db.C(COLLECTION).Find(bson.M{}).All(&movies)
	fmt.Println(err)
	return movies
}

// GetUser ...
func (s *Storage) GetUser(idsub int) (listing.User, error) {
	var id = string(idsub)
	var user listing.User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

// AddUser ...
func (s *Storage) AddUser(user authenticating.User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

// Delete ...
func (s *Storage) Delete(user User) error {
	err := db.C(COLLECTION).Remove(&user)
	return err
}

// Update ...
func (s *Storage) Update(user User) error {
	err := db.C(COLLECTION).UpdateId(user.UID, &user)
	return err
}

// TODO: to be implemented

// JWTData is a struct with the structure of the jwt data
type JWTData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

// SignUp ...
func (s *Storage) SignUp(username string, password string, firstname string, lastname string) (interface{}, error) {
	returnObjectMap = make(map[string]interface{})
	var user authenticating.User
	var result []byte

	user.UserName = username
	user.Password = password
	user.FirstName = firstname
	user.LastName = lastname
	user.Gender = ""

	err := db.C(COLLECTION).Insert(&user)
	// fmt.Println("user: ", user)
	if err != nil {
		return nil, err
	}

	// err = db.C(COLLECTION).Find(bson.M{"username": user.UserName}).One(&user)

	// if err != nil {
	// 	return nil, err
	// }

	// var userUID = user.UID.Hex()

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

	fmt.Println("got here")

	returnObjectMap["token"] = result
	returnObjectMap["userUID"] = ""

	return returnObjectMap, nil

}

// Login ...
func (s *Storage) Login(username string, password string) (interface{}, error) {
	returnObjectMap = make(map[string]interface{})
	var user listing.User
	var result []byte

	user.UserName = username
	user.Password = password

	err := db.C(COLLECTION).Find(bson.M{"username": user.UserName, "password": user.Password}).One(&user)

	if err != nil {
		return nil, err
	}

	var userUID = user.UID.Hex()
	// Demo - in case no db
	// if user.UserName == "admin" && user.Password == "password" {

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},

		CustomClaims: map[string]string{
			"userid": userUID,
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

	returnObjectMap["token"] = result
	returnObjectMap["userUID"] = userUID

	return returnObjectMap, nil
}

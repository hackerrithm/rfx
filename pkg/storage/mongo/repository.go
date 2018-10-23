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

type MoviesDAO struct {
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "user"
)

func NewStorage() (*MoviesDAO, error) {
	session, err := mgo.Dial("mongodb://localhost/test1")
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB("test1")
	fmt.Println("connected")
	return &MoviesDAO{db.Name}, err
}

// GetAllUsers ...
func (m *MoviesDAO) GetAllUsers() []listing.User {
	var movies []listing.User
	err := db.C(COLLECTION).Find(bson.M{}).All(&movies)
	fmt.Println(err)
	return movies
}

// GetUser ...
func (m *MoviesDAO) GetUser(idsub int) (listing.User, error) {
	var id = string(idsub)
	var movie listing.User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&movie)
	return movie, err
}

// AddUser ...
func (m *MoviesDAO) AddUser(movie authenticating.User) error {
	err := db.C(COLLECTION).Insert(&movie)
	return err
}

// Delete ...
func (m *MoviesDAO) Delete(movie User) error {
	err := db.C(COLLECTION).Remove(&movie)
	return err
}

// Update ...
func (m *MoviesDAO) Update(movie User) error {
	err := db.C(COLLECTION).UpdateId(movie.UID, &movie)
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

// Login ...
func (m *MoviesDAO) Login(username string, password string) ([]byte, error) {

	var user listing.User
	var result []byte

	user.UserName = username
	user.Password = password

	err := db.C(COLLECTION).Find(bson.M{"username": user.UserName, "password": user.Password}).One(&user)

	if err != nil {
		return nil, err
	}

	// Demo - in real case scenario you'd check this against your database
	// if user.UserName == "admin" && user.Password == "password" {

	fmt.Println("we aight")
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

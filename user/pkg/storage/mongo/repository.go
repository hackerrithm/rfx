package mongo

import (
	"encoding/json"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hackerrithm/longterm/rfx/user/pkg/authenticating"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// SECRET ...
	SECRET = "12This98Is34A76String56Used65As78Secret01"
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
	user.SetPassword(password)
	user.FirstName = firstname
	user.LastName = lastname

	err := db.C(COLLECTION).Insert(&user)
	if err != nil {
		log.Printf("%s", err)
	}

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},

		CustomClaims: map[string]string{
			"userid": "u1",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRET))
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

// Login ...
func (s *Storage) Login(username string, password string) (map[string]interface{}, error) {
	returnObjectMap = make(map[string]interface{})
	var user authenticating.User
	user.UserName = username

	err := db.C(COLLECTION).Find(bson.M{"username": user.UserName}).One(&user)

	if err != nil {
		return nil, err
	}

	if !user.IsCredentialsVerified(password, user.Password) {
		log.Println("error")
		return nil, err
	}

	var userUID = user.UID.Hex()

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},

		CustomClaims: map[string]string{
			"userid": "u1",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRET))
	if err != nil {
		log.Println("StatusUnauthorized ", err)
	}

	returnObjectMap["token"] = tokenString
	returnObjectMap["userUID"] = userUID

	return returnObjectMap, nil
}

func getByteToken(token *jwt.Token) (interface{}, error) {
	if jwt.SigningMethodHS256 != token.Method {
		log.Println("Invalid signing algorithm")
	}
	return []byte(SECRET), nil
}

func parseWithClaims(jwtToken string) (*jwt.Token, error) {
	cl, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, getByteToken)
	if err != nil {
		log.Println("error in parseWithClaims")
	}
	return cl, nil
}

// Profile ...
func (s *Storage) Profile(jwtToken string, UUID string) ([]byte, error) {
	claims, err := parseWithClaims(jwtToken)
	if err != nil {
		log.Println(err)
	}

	data := claims.Claims.(*JWTData)

	userID := data.CustomClaims["userid"]
	log.Println("claim ", userID)

	// fetch some data based on the userID and then send that data back to the user in JSON format
	user, err := GetUser(UUID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser ...
func GetUser(ID string) ([]byte, error) {
	var user authenticating.User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(ID)).One(&user)
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	return []byte(json), err
}

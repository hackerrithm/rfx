package mongodb

import (
	"encoding/json"
	"log"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hackerrithm/longterm/rfx/user/domain"
	"github.com/hackerrithm/longterm/rfx/user/engine"
)

const (
	// SECRET ...
	SECRET = "12This98Is34A76String56Used65As78Secret01"
)

var returnObjectMap map[string]interface{}

// TODO: to be implemented

// JWTData is a struct with the structure of the jwt data
type JWTData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

type (
	userRepository struct {
		session *mgo.Session
	}
)

var (
	userCollection = "user"
)

func newUserRepository(session *mgo.Session) engine.UserRepository {
	return &userRepository{session}
}

func (r userRepository) Put(c context.Context, u *domain.User) map[string]interface{} {
	s := r.session.Clone()
	defer s.Close()

	returnObjectMap = make(map[string]interface{})
	var user domain.User
	var result []byte

	user.UserName = u.UserName
	user.SetPassword(u.Password)
	user.FirstName = u.FirstName
	user.LastName = u.LastName

	col := s.DB("test1").C(userCollection)
	if u.ID == 0 {
		u.ID = getNextSequence(s, userCollection)
	}

	// col.Upsert(bson.M{"_id": u.ID}, u)
	col.Insert(&user)

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

	return returnObjectMap
}

func (r userRepository) List(c context.Context, query *engine.Query) []*domain.User {
	s := r.session.Clone()
	defer s.Close()

	col := s.DB("test1").C(userCollection)
	g := []*domain.User{}
	q := translateQuery(col, query)
	q.All(&g)

	return g
}

func (r userRepository) Read(c context.Context, username, password string) map[string]interface{} {
	s := r.session.Clone()
	defer s.Close()

	returnObjectMap = make(map[string]interface{})
	var user *domain.User

	err := s.DB("test1").C(userCollection).Find(bson.M{"username": username}).One(&user)

	if err != nil {
		return nil
	}

	if !user.IsCredentialsVerified(password, user.Password) {
		log.Println("error")
		return nil
	}

	var userUID = user.ID

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},

		CustomClaims: map[string]string{
			"userid": string(userUID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRET))
	if err != nil {
		log.Println("StatusUnauthorized ", err)
	}

	returnObjectMap["token"] = tokenString
	returnObjectMap["userUID"] = userUID

	return returnObjectMap
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
func (r userRepository) Profile(c context.Context, jwtToken string, ID string) []byte {
	claims, err := parseWithClaims(jwtToken)
	if err != nil {
		log.Println(err)
	}

	data := claims.Claims.(*JWTData)

	userID := data.CustomClaims["userid"]
	log.Println("claim ", userID)

	s := r.session.Clone()
	defer s.Close()

	var user domain.User
	err = s.DB("test1").C(userCollection).FindId(bson.ObjectIdHex(ID)).One(&user)
	if err != nil {
		return nil
	}

	json, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return nil
	}

	return []byte(json)
}

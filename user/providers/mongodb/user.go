package mongodb

import (
	"log"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/hackerrithm/longterm/rfx/user/domain"
	"github.com/hackerrithm/longterm/rfx/user/engine"
)

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

func (r userRepository) Insert(c context.Context, u *domain.User) string {
	s := r.session.Clone()
	defer s.Close()

	var user domain.User

	user.UserName = u.UserName
	user.SetPassword(u.Password)
	user.FirstName = u.FirstName
	user.LastName = u.LastName

	col := s.DB("test1").C(userCollection)
	err := col.Insert(&user)
	if err != nil {
		log.Println(err)
	}

	return "ok"
}

func (r userRepository) Put(c context.Context, u *domain.User) string {
	s := r.session.Clone()
	defer s.Close()

	var user domain.User

	user.UserName = u.UserName
	user.SetPassword(u.Password)
	user.FirstName = u.FirstName
	user.LastName = u.LastName

	col := s.DB("test1").C(userCollection)
	col.Upsert(bson.M{"_id": u.ID}, u)
	return "ok"
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

func (r userRepository) Read(c context.Context, username, password string) *domain.User {
	s := r.session.Clone()
	defer s.Close()

	var user *domain.User

	err := s.DB("test1").C(userCollection).Find(bson.M{"username": username}).One(&user)

	if err != nil {
		return nil
	}

	if !user.IsCredentialsVerified(password, user.Password) {
		log.Println("error")
		return nil
	}

	return user
}

// Profile ...
func (r userRepository) Profile(c context.Context, jwtToken string, ID string) *domain.User {
	s := r.session.Clone()
	defer s.Close()

	var user domain.User
	err := s.DB("test1").C(userCollection).FindId(bson.ObjectIdHex(ID)).One(&user)
	if err != nil {
		return nil
	}
	return &user
}

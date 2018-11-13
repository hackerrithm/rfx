package mongodb

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/hackerrithm/longterm/rfx/user/engine"
)

type (
	storageFactory struct {
		session *mgo.Session
	}
)

// NewStorage creates a new instance of this mongodb storage factory
func NewStorage(url string) engine.StorageFactory {
	session, _ := mgo.DialWithTimeout(url, 10*time.Second)
	ensureIndexes(session)
	return &storageFactory{session}
}

// NewUserRepository creates a new datastore user repository
func (f *storageFactory) NewUserRepository() engine.UserRepository {
	return newUserRepository(f.session)
}

func ensureIndexes(s *mgo.Session) {
	index := mgo.Index{
		Key:        []string{"date"},
		Background: true,
	}
	c := s.DB("test1").C(userCollection)
	c.EnsureIndex(index)
}

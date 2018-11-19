package mongodb

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/hackerrithm/longterm/rfx/user/engine"
)

type (
	storageFactory struct {
		session *mgo.Session
	}
)

var (
	userRepositoryInstance engine.UserRepository
	userRepositoryOnce     sync.Once
)

// NewStorage creates a new instance of this mongodb storage factory
func NewStorage(url string) engine.StorageFactory {
	session, _ := mgo.DialWithTimeout(url, 10*time.Second)
	ensureIndexes(session)
	return &storageFactory{session}
}

// NewUserRepository creates a new datastore user repository
func (f *storageFactory) NewUserRepository() engine.UserRepository {
	userRepositoryOnce.Do(func() {
		userRepositoryInstance = newUserRepository(f.session)
	})
	return userRepositoryInstance
}

func ensureIndexes(s *mgo.Session) {
	index := mgo.Index{
		Key:        []string{"date"},
		Background: true,
	}
	c := s.DB("test1").C(userCollection)
	c.EnsureIndex(index)
}

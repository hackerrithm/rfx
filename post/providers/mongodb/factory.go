package mongodb

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/hackerrithm/longterm/rfx/post/engine"
)

type (
	storageFactory struct {
		session *mgo.Session
	}
)

var (
	postRepositoryInstance engine.PostRepository
	postRepositoryOnce     sync.Once
)

// NewStorage creates a new instance of this mongodb storage factory
func NewStorage(url string) engine.StorageFactory {
	session, _ := mgo.DialWithTimeout(url, 10*time.Second)
	ensureIndexes(session)
	return &storageFactory{session}
}

// NewPostRepository creates a new datastore Post repository
func (f *storageFactory) NewPostRepository() engine.PostRepository {
	postRepositoryOnce.Do(func() {
		postRepositoryInstance = newPostRepository(f.session)
	})
	return postRepositoryInstance
}

func ensureIndexes(s *mgo.Session) {
	index := mgo.Index{
		Key:        []string{"date"},
		Background: true,
	}
	c := s.DB("test1").C(postCollection)
	c.EnsureIndex(index)
}

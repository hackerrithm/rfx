package mongodb

import (
	"time"

	"github.com/hackerrithm/longterm/rfx/task/engine"
	mgo "gopkg.in/mgo.v2"
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

// NewTaskRepository creates a new datastore task repository
func (f *storageFactory) NewTaskRepository() engine.TaskRepository {
	return newTaskRepository(f.session)
}

func ensureIndexes(s *mgo.Session) {
	index := mgo.Index{
		Key:        []string{"date"},
		Background: true,
	}
	c := s.DB("").C(taskCollection)
	c.EnsureIndex(index)
}

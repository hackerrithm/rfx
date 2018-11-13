package mongodb

import (
	"context"

	"github.com/hackerrithm/longterm/rfx/task/domain"
	"github.com/hackerrithm/longterm/rfx/task/engine"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	taskRepository struct {
		session *mgo.Session
	}
)

var (
	taskCollection = "task"
)

func newTaskRepository(session *mgo.Session) engine.TaskRepository {
	return &taskRepository{session}
}

func (r taskRepository) Put(c context.Context, g *domain.Task) {
	s := r.session.Clone()
	defer s.Close()

	col := s.DB("test1").C(taskCollection)
	if g.ID == 0 {
		g.ID = getNextSequence(s, taskCollection)
	}
	col.Upsert(bson.M{"_id": g.ID}, g)
}

func (r taskRepository) List(c context.Context, query *engine.Query) []*domain.Task {
	s := r.session.Clone()
	defer s.Close()

	col := s.DB("test1").C(taskCollection)
	g := []*domain.Task{}
	q := translateQuery(col, query)
	q.All(&g)

	return g
}

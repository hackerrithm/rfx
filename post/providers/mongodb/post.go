package mongodb

import (
	"log"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/hackerrithm/longterm/rfx/post/domain"
	"github.com/hackerrithm/longterm/rfx/post/engine"
)

type (
	postRepository struct {
		session *mgo.Session
	}
)

var (
	postCollection = "post"
)

func newPostRepository(session *mgo.Session) engine.PostRepository {
	return &postRepository{session}
}

func (r postRepository) Insert(c context.Context, p *domain.Post) string {
	s := r.session.Clone()
	defer s.Close()

	var post domain.Post

	post.Author = p.Author
	post.Category = p.Category
	post.Topic = p.Topic
	post.ContentPhoto = p.ContentPhoto
	post.ContentText = p.ContentText

	col := s.DB("test1").C(postCollection)
	err := col.Insert(&post)
	if err != nil {
		log.Println(err)
	}

	return "ok"
}

func (r postRepository) Put(c context.Context, p *domain.Post) string {
	s := r.session.Clone()
	defer s.Close()

	var post domain.Post

	post.Author = p.Author
	post.Category = p.Category
	post.Topic = p.Topic
	post.ContentPhoto = p.ContentPhoto
	post.ContentText = p.ContentText

	col := s.DB("test1").C(postCollection)
	col.Upsert(bson.M{"_id": p.ID}, p)
	return "ok"
}

func (r postRepository) List(c context.Context, query *engine.Query) []*domain.Post {
	s := r.session.Clone()
	defer s.Close()

	col := s.DB("test1").C(postCollection)
	p := []*domain.Post{}
	q := translateQuery(col, query)
	q.All(&p)

	return p
}

func (r postRepository) Read(c context.Context, id string) *domain.Post {
	s := r.session.Clone()
	defer s.Close()

	var post *domain.Post

	err := s.DB("test1").C(postCollection).Find(bson.M{"_id": id}).One(&post)

	if err != nil {
		return nil
	}

	return post
}

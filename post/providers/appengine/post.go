// +build appengine

package appengine

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/hackerrithm/longterm/rfx/post/domain"
	"github.com/hackerrithm/longterm/rfx/post/engine"
)

type (
	postRepository struct{}

	// Post is the internal struct we use for persistence
	// it allows us to attach the datastore PropertyLoadSaver
	// methods to the struct that we don't own
	post struct {
		domain.Post
	}
)

var (
	postKind = "post"
)

func newPostRepository() engine.PostRepository {
	return &postRepository{}
}

func (r postRepository) Put(c context.Context, g *domain.Post) {
	d := &post{*g}
	k := datastore.NewIncompleteKey(c, postKind, guestbookEntityKey(c))
	datastore.Put(c, k, d)
}

func (r postRepository) List(c context.Context, query *engine.Query) []*domain.Post {
	g := []*post{}
	q := translateQuery(postKind, query)
	q = q.Ancestor(guestbookEntityKey(c))

	k, _ := q.GetAll(c, &g)
	o := make([]*domain.Post, len(g))
	for i := range g {
		o[i] = &g[i].Post
		o[i].ID = k[i].IntID()
	}
	return o
}

func guestbookEntityKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "guestbook", "", 1, nil)
}

func (x *Post) Load(props []datastore.Property) error {
	for _, prop := range props {
		switch prop.Name {
		case "author":
			x.Author = prop.Value.(string)
		case "content":
			x.Content = prop.Value.(string)
		case "date":
			x.Date = prop.Value.(time.Time)
		}
	}
	return nil
}

func (x *Post) Save() ([]datastore.Property, error) {
	ps := []datastore.Property{
		datastore.Property{Name: "author", Value: x.Author, NoIndex: true, Multiple: false},
		datastore.Property{Name: "content", Value: x.Content, NoIndex: true, Multiple: false},
		datastore.Property{Name: "date", Value: x.Date, NoIndex: false, Multiple: false},
	}
	return ps, nil
}

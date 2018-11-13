// +build appengine

package appengine

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/hackerrithm/longterm/rfx/user/domain"
	"github.com/hackerrithm/longterm/rfx/user/engine"
)

type (
	userRepository struct{}

	// user is the internal struct we use for persistence
	// it allows us to attach the datastore PropertyLoadSaver
	// methods to the struct that we don't own
	user struct {
		domain.User
	}
)

var (
	userKind = "user"
)

func newUserRepository() engine.UserRepository {
	return &userRepository{}
}

func (r userRepository) Put(c context.Context, g *domain.User) {
	d := &user{*g}
	k := datastore.NewIncompleteKey(c, userKind, guestbookEntityKey(c))
	datastore.Put(c, k, d)
}

func (r userRepository) List(c context.Context, query *engine.Query) []*domain.User {
	g := []*user{}
	q := translateQuery(userKind, query)
	q = q.Ancestor(guestbookEntityKey(c))

	k, _ := q.GetAll(c, &g)
	o := make([]*domain.User, len(g))
	for i := range g {
		o[i] = &g[i].User
		o[i].ID = k[i].IntID()
	}
	return o
}

func guestbookEntityKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "guestbook", "", 1, nil)
}

func (x *user) Load(props []datastore.Property) error {
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

func (x *user) Save() ([]datastore.Property, error) {
	ps := []datastore.Property{
		datastore.Property{Name: "author", Value: x.Author, NoIndex: true, Multiple: false},
		datastore.Property{Name: "content", Value: x.Content, NoIndex: true, Multiple: false},
		datastore.Property{Name: "date", Value: x.Date, NoIndex: false, Multiple: false},
	}
	return ps, nil
}

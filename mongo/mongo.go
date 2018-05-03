package mongo

import (
	"github.com/gin-contrib/sessions"
	gsessions "github.com/gorilla/sessions"
	"github.com/kidstuff/mongostore"
	mgo "gopkg.in/mgo.v2"
)

type Store interface {
	sessions.Store
}

func NewStore(c *mgo.Collection, maxAge int, ensureTTL bool, keyPairs ...[]byte) Store {
	return &store{mongostore.NewMongoStore(c, maxAge, ensureTTL, keyPairs...)}
}

type store struct {
	*mongostore.MongoStore
}

func (c *store) Options(options sessions.Options) {
	c.MongoStore.Options = &gsessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

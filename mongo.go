package sessions

import (
	"github.com/gorilla/sessions"
	"gopkg.in/go-playground/mongostore.v4"
	mgo "gopkg.in/mgo.v2"
)

type MongoStore interface {
	Store
}

func NewMongoStore(s *mgo.Session, collectionName string, options *sessions.Options, ensureTTL bool, keyPairs ...[]byte) MongoStore {
	store := mongostore.New(s, collectionName, options, ensureTTL, keyPairs...)

	return &mongoStore{store}
}

type mongoStore struct {
	*mongostore.MongoStore
}

func (c *mongoStore) Options(options Options) {
	c.MongoStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

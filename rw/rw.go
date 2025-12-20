package rw

import (
	"github.com/gin-contrib/sessions"
	gsessions "github.com/gorilla/sessions"
	"net/http"
)

type Store interface {
	sessions.Store
}

// NewStore create a session store split read and write channel
func NewStore(read, write sessions.Store) Store {
	return &store{read, write}
}

type store struct {
	read  sessions.Store
	write sessions.Store
}

// Get should return a cached session.
func (c *store) Get(r *http.Request, name string) (*gsessions.Session, error) {
	return c.read.Get(r, name)
}

// New should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (c *store) New(r *http.Request, name string) (*gsessions.Session, error) {
	return c.write.New(r, name)
}

// Save should persist session to the underlying store implementation.
func (c *store) Save(r *http.Request, w http.ResponseWriter, s *gsessions.Session) error {
	return c.write.Save(r, w, s)
}

func (c *store) Options(options sessions.Options) {
	c.read.Options(options)
	c.write.Options(options)
}
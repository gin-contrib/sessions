package sessions

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

const (
	DefaultKey  = "github.com/gin-contrib/sessions"
	errorFormat = "[sessions] ERROR! %s\n"
)

type Store interface {
	sessions.Store
	Options(Options)
}

// Options stores configuration for a session or session store.
// Fields are a subset of http.Cookie fields.
type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

// Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type Session interface {
	// Get returns the session value associated to the given key.
	Get(key interface{}) (interface{}, error)
	// Set sets the session value associated to the given key.
	Set(key interface{}, val interface{}) error
	// Delete removes the session value associated to the given key.
	Delete(key interface{}) error
	// Clear deletes all values in the session.
	Clear() error
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string) error
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) ([]interface{}, error)
	// Options sets confuguration for a session.
	Options(Options) error
	// Save saves all sessions used during the current request.
	Save() error
}

func Sessions(name string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &session{name, c.Request, store, nil, false, c.Writer}
		c.Set(DefaultKey, s)
		defer context.Clear(c.Request)
		c.Next()
	}
}

type session struct {
	name    string
	request *http.Request
	store   Store
	session *sessions.Session
	written bool
	writer  http.ResponseWriter
}

func (s *session) Get(key interface{}) (interface{}, error) {
	v, err := s.Session()
	if err != nil{
		return nil, err
	}
	return v.Values[key], nil
}

func (s *session) Set(key interface{}, val interface{}) error {
	v, err := s.Session()
	if err != nil{
		return err
	}
	v.Values[key] = val
	s.written = true
	return nil
}

func (s *session) Delete(key interface{}) error {
	v, err := s.Session()
	if err != nil{
		return err
	}
	delete(v.Values, key)
	s.written = true
	return nil
}

func (s *session) Clear() error {
	v, err := s.Session()
	if err != nil{
		return err
	}
	for key := range v.Values {
		s.Delete(key)
	}
	return nil
}

func (s *session) AddFlash(value interface{}, vars ...string) error {
	v, err := s.Session()
	if err != nil{
		return err
	}
	v.AddFlash(value, vars...)
	s.written = true
	return nil
}

func (s *session) Flashes(vars ...string) ([]interface{}, error) {
	s.written = true
	v, err := s.Session()
	if err != nil{
		return nil, err
	}
	return v.Flashes(vars...), nil
}

func (s *session) Options(options Options) error {
	v, err := s.Session()
	if err != nil{
		return err
	}
	v.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
	return nil
}

func (s *session) Save() error {
	if s.Written() {
		v, err := s.Session()
		if err != nil{
			return err
		}
		e := v.Save(s.request, s.writer)
		if e == nil {
			s.written = false
		}
		return e
	}
	return nil
}

func (s *session) Session() (*sessions.Session, error) {
	var err error
	if s.session == nil {
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session, err
}

func (s *session) Written() bool {
	return s.written
}

// shortcut to get session
func Default(c *gin.Context) Session {
	return c.MustGet(DefaultKey).(Session)
}

package cassandra

import (
	"context"
	"encoding/base32"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/gocql/gocql"
	"github.com/gorilla/securecookie"
	gsessions "github.com/gorilla/sessions"
)

const (
	defaultTable = "sessions"
	defaultTTL   = 86400
)

// Store implements sessions.Store using Apache Cassandra.
type Store struct {
	session *gocql.Session
	table   string
	codecs  []securecookie.Codec
	options *gsessions.Options
}

// NewStore creates a Cassandra session store.
//
// maxAge is the session lifetime in seconds.
//
// Example:
//
//	cluster := gocql.NewCluster("127.0.0.1")
//	cluster.Keyspace = "sessions"
//
//	store, err := NewStore(
//		cluster,
//		3600,
//		[]byte("authentication-key"),
//		[]byte("encryption-key"),
//	)
func NewStore(
	cluster *gocql.ClusterConfig,
	maxAge int,
	keyPairs ...[]byte,
) (*Store, error) {
	if cluster == nil {
		return nil, errors.New("cassandra: cluster config is nil")
	}

	if len(keyPairs) == 0 {
		return nil, errors.New("cassandra: at least one key pair is required")
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("cassandra: create session: %w", err)
	}

	store := &Store{
		session: session,
		table:   defaultTable,
		codecs:  securecookie.CodecsFromPairs(keyPairs...),
		options: &gsessions.Options{
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		},
	}

	return store, nil
}

// NewStoreWithOptions creates a Cassandra session store with custom options.
func NewStoreWithOptions(
	cluster *gocql.ClusterConfig,
	options *gsessions.Options,
	keyPairs ...[]byte,
) (*Store, error) {
	if cluster == nil {
		return nil, errors.New("cassandra: cluster config is nil")
	}

	if len(keyPairs) == 0 {
		return nil, errors.New("cassandra: at least one key pair is required")
	}

	if options == nil {
		options = &gsessions.Options{
			Path:     "/",
			MaxAge:   defaultTTL,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("cassandra: create session: %w", err)
	}

	return &Store{
		session: session,
		table:   defaultTable,
		codecs:  securecookie.CodecsFromPairs(keyPairs...),
		options: options,
	}, nil
}

// Close closes the Cassandra session.
func (s *Store) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

// Table changes the Cassandra table used for sessions.
func (s *Store) Table(table string) *Store {
	if table != "" {
		s.table = table
	}

	return s
}

// Options sets the session options, satisfying sessions.Store.
func (s *Store) Options(options sessions.Options) {
	s.options = options.ToGorillaOptions()
}

// New returns a new session.
func (s *Store) New(
	request *http.Request,
	name string,
) (*gsessions.Session, error) {
	session := gsessions.NewSession(s, name)

	session.Options = cloneOptions(s.options)
	session.IsNew = true

	cookie, err := request.Cookie(name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return session, nil
		}

		return session, err
	}

	if cookie.Value == "" {
		return session, nil
	}

	var id string

	if err := securecookie.DecodeMulti(
		name,
		cookie.Value,
		&id,
		s.codecs...,
	); err != nil {
		// A tampered or stale cookie behaves as a new session.
		return session, nil
	}

	ok, err := s.load(request.Context(), id, session)
	if err != nil {
		return session, err
	}

	if !ok {
		// The session is missing or expired in Cassandra.
		return session, nil
	}

	session.ID = id
	session.IsNew = false

	return session, nil
}

// load reads and decodes a session from Cassandra.
//
// It returns false when the session no longer exists (missing or expired)
// or when the stored data cannot be decoded, so the caller can fall back to
// a fresh session.
func (s *Store) load(
	ctx context.Context,
	id string,
	session *gsessions.Session,
) (bool, error) {
	query := fmt.Sprintf(
		`SELECT data FROM %s WHERE session_id = ?`,
		s.table,
	)

	var data []byte

	if err := s.session.Query(query, id).
		WithContext(ctx).
		Scan(&data); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("cassandra: load session: %w", err)
	}

	values := make(map[interface{}]interface{})

	if err := securecookie.DecodeMulti(
		session.Name(),
		string(data),
		&values,
		s.codecs...,
	); err != nil {
		// Corrupt data should behave as a new session.
		return false, nil
	}

	session.Values = values

	return true, nil
}

// Get implements sessions.Store.
func (s *Store) Get(
	request *http.Request,
	name string,
) (*gsessions.Session, error) {
	return s.New(request, name)
}

// Save persists the session.
func (s *Store) Save(
	request *http.Request,
	response http.ResponseWriter,
	session *gsessions.Session,
) error {
	if session == nil {
		return errors.New("cassandra: session is nil")
	}

	// MaxAge < 0 means delete the session.
	if session.Options.MaxAge < 0 {
		if session.ID != "" {
			if err := s.delete(request.Context(), session.ID); err != nil {
				return err
			}
		}

		http.SetCookie(
			response,
			newCookie(session.Name(), "", session.Options),
		)

		return nil
	}

	if session.ID == "" {
		id, err := generateID()
		if err != nil {
			return err
		}

		session.ID = id
	}

	encoded, err := securecookie.EncodeMulti(
		session.Name(),
		session.Values,
		s.codecs...,
	)
	if err != nil {
		return fmt.Errorf("cassandra: encode session: %w", err)
	}

	ttl := session.Options.MaxAge

	if ttl <= 0 {
		ttl = defaultTTL
	}

	expiresAt := time.Now().Add(
		time.Duration(ttl) * time.Second,
	)

	query := fmt.Sprintf(
		`INSERT INTO %s
		 (session_id, data, expires_at)
		 VALUES (?, ?, ?)
		 USING TTL ?`,
		s.table,
	)

	if err := s.session.Query(
		query,
		session.ID,
		encoded,
		expiresAt,
		ttl,
	).
		WithContext(request.Context()).
		Exec(); err != nil {
		return fmt.Errorf("cassandra: save session: %w", err)
	}

	cookieValue, err := securecookie.EncodeMulti(
		session.Name(),
		session.ID,
		s.codecs...,
	)
	if err != nil {
		return fmt.Errorf("cassandra: encode session id: %w", err)
	}

	http.SetCookie(
		response,
		newCookie(session.Name(), cookieValue, session.Options),
	)

	session.IsNew = false

	return nil
}

// delete removes a session from Cassandra.
func (s *Store) delete(
	ctx context.Context,
	sessionID string,
) error {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE session_id = ?`,
		s.table,
	)

	if err := s.session.Query(query, sessionID).
		WithContext(ctx).
		Exec(); err != nil {
		return fmt.Errorf("cassandra: delete session: %w", err)
	}

	return nil
}

// Clear removes all sessions.
//
// This should generally only be used in tests or administrative operations.
func (s *Store) Clear(ctx context.Context) error {
	query := fmt.Sprintf(
		`TRUNCATE %s`,
		s.table,
	)

	if err := s.session.Query(query).
		WithContext(ctx).
		Exec(); err != nil {
		return fmt.Errorf("cassandra: clear sessions: %w", err)
	}

	return nil
}

// generateID returns a cryptographically random, URL-safe session identifier.
func generateID() (string, error) {
	key := securecookie.GenerateRandomKey(32)
	if key == nil {
		return "", errors.New("cassandra: generate session id")
	}

	return strings.TrimRight(
		base32.StdEncoding.EncodeToString(key),
		"=",
	), nil
}

func cloneOptions(options *gsessions.Options) *gsessions.Options {
	if options == nil {
		return &gsessions.Options{
			Path:     "/",
			MaxAge:   defaultTTL,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
	}

	copy := *options

	return &copy
}

func newCookie(
	name string,
	value string,
	options *gsessions.Options,
) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: options.SameSite,
	}

	if options.MaxAge < 0 {
		cookie.Expires = time.Unix(1, 0)
	}

	return cookie
}

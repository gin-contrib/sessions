package postgres

import (
	"database/sql"

	"github.com/antonlindstrom/pgstore"
	"github.com/gin-contrib/sessions"
	gsessions "github.com/gorilla/sessions"
)

type Store interface {
	sessions.Store
}

func NewStore(url string, keyPairs ...[]byte) (Store, error) {
	x, err := pgstore.NewPGStore(url, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{x}, nil
}

func NewStoreFromPool(db *sql.DB, keyPairs ...[]byte) (Store, error) {
	x, err := pgstore.NewPGStoreFromPool(db, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{x}, nil
}

type store struct {
	*pgstore.PGStore
}

func (s *store) Options(options sessions.Options) {
	s.PGStore.Options = &gsessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: options.SameSite,
	}
}

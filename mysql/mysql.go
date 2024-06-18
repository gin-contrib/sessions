package mysql

import (
	"database/sql"

	"github.com/weisskopfjens/mysqlstore"
	"github.com/weisskopfjens/sessions"
)

type Store interface {
	sessions.Store
}

type store struct {
	*mysqlstore.MySQLStore
}

var _ Store = new(store)

func NewStore(db *sql.DB, keyPairs ...[]byte) (Store, error) {
	p, err := mysqlstore.NewMySQLStoreFromConnection(db, "sessions", "/", 3600, keyPairs...)
	if err != nil {
		return nil, err
	}

	return &store{p}, nil
}

func (s *store) Options(options sessions.Options) {
	s.MySQLStore.Options = options.ToGorillaOptions()
}

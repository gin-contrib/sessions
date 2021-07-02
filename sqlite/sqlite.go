package sqlite

import (
	"crawshaw.io/sqlite/sqlitex"
	"github.com/gin-contrib/sessions"
	"github.com/terem42/sqlite3store"
)

type Store interface {
	sessions.Store
}

func NewStore(db_file_location string, db_pool_size int, db_table_name string, sessions_options *sessions.Options, keyPairs ...[]byte) (Store, error) {
	s, err := sqlite3store.NewSqliteStore(db_file_location, db_pool_size, db_table_name, sessions_options.ToGorillaOptions(), keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

//NewSqliteStoreFromExistingPool(dbpool *sqlitex.Pool, tableName string, sessions_options sessions.Options, keyPairs ...[]byte) (*SqliteStore, error) {
func NewStoreFromExistingPool(dbpool *sqlitex.Pool, tableName string, sessions_options *sessions.Options, keyPairs ...[]byte) (Store, error) {
	s, err := sqlite3store.NewSqliteStoreFromExistingPool(dbpool, tableName, sessions_options.ToGorillaOptions(), keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

type store struct {
	*sqlite3store.SqliteStore
}

func (c *store) Options(options sessions.Options) {
	c.SqliteStore.Options = options.ToGorillaOptions()
}

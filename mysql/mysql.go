package mysql

import (
	"github.com/gin-contrib/sessions"
	"github.com/srinathgs/mysqlstore"
	"time"
)

type SessionStore interface {
	sessions.Store
}

type store struct {
	*mysqlstore.MySQLStore
}

func NewStore(endpoint string, tableName string, path string, maxAge int, keyPairs ...[]byte) (store, error) {
	p, err := mysqlstore.NewMySQLStore(endpoint, tableName, path, maxAge, keyPairs...)
	if err != nil {
		return store{},err
	}

	return store{p}, nil
}

func (s *store) Options(options sessions.Options) {
	s.MySQLStore.Options = options.ToGorillaOptions()
}

func (s *store) Close() {
	s.MySQLStore.Close()
}

func (s *store) Cleanup(interval time.Duration) (chan<- struct{}, <-chan struct{}){
	n,m := s.MySQLStore.Cleanup(interval)
	return n, m
}

func (s *store) StopCleanup(quit chan<- struct{}, done <-chan struct{}) {
	s.MySQLStore.StopCleanup(quit,done)
}
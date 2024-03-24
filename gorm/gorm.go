package gorm

import (
	"time"

	"github.com/wader/gormstore/v2"
	"github.com/weisskopfjens/sessions"
	"gorm.io/gorm"
)

type Store interface {
	sessions.Store
}

func NewStore(d *gorm.DB, expiredSessionCleanup bool, keyPairs ...[]byte) Store {
	s := gormstore.New(d, keyPairs...)
	if expiredSessionCleanup {
		quit := make(chan struct{})
		go s.PeriodicCleanup(1*time.Hour, quit)
	}
	return &store{s}
}

type store struct {
	*gormstore.Store
}

func (s *store) Options(options sessions.Options) {
	s.Store.SessionOpts = options.ToGorillaOptions()
}

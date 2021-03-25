package gorm

import (
	"github.com/gin-contrib/sessions"
	"github.com/wader/gormstore/v2"
	"gorm.io/gorm"
)

type Store interface {
	sessions.Store
}

func NewStore(d *gorm.DB, keyPairs ...[]byte) Store {
	return &store{gormstore.New(d, keyPairs...)}
}

type store struct {
	*gormstore.Store
}

func (s *store) Options(options sessions.Options) {
	s.Store.SessionOpts = options.ToGorillaOptions()
}

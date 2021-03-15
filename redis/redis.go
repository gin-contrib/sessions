package redis

import (
	"context"
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/go-redis/redis/v8"
	"github.com/rbcervilla/redisstore/v8"
)

type Store interface {
	sessions.Store
}

// NewStore - create new session store with given redis client interface
func NewStore(context context.Context, client redis.UniversalClient) (sessions.Store, error) {
	innerStore, err := redisstore.NewRedisStore(context, client)
	if err != nil {
		return nil, err
	}
	return &store{innerStore}, nil
}

type store struct {
	*redisstore.RedisStore
}

// GetRedisStore get the actual working store.
// Ref: https://godoc.org/github.com/boj/redistore#RediStore
func GetRedisStore(s Store) (err error, rediStore *redisstore.RedisStore) {
	realStore, ok := s.(*store)
	if !ok {
		err = errors.New("unable to get the redis store: Store isn't *store")
		return
	}

	rediStore = realStore.RedisStore
	return
}

// SetKeyPrefix sets the key prefix in the redis database.
func SetKeyPrefix(s Store, prefix string) error {
	err, rediStore := GetRedisStore(s)
	if err != nil {
		return err
	}

	rediStore.KeyPrefix(prefix)
	return nil
}

func (c *store) Options(options sessions.Options) {
	c.RedisStore.Options(*options.ToGorillaOptions())
}

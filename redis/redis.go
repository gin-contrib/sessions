package redis

import (
	"github.com/gin-contrib/sessions"
	"github.com/joelrose/redisstore"
	"github.com/joelrose/redisstore/adapter"
	"github.com/redis/go-redis/v9"
)

type Store interface {
	sessions.Store
}

// Client is a redis client that implements the redisstore.Client interface.
//
// keyPrefix is the prefix for all keys in redis.
//
// keyPairs are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewStore(client *redis.Client, keyPrefix string, keyPairs ...[]byte) Store {
	s := redisstore.New(
		adapter.WithGoRedis(client),
		keyPairs,
		redisstore.WithKeyPrefix(keyPrefix),
	)

	return &store{s}
}

type store struct {
	*redisstore.Store
}

func (c *store) Options(options sessions.Options) {
	c.Store.SetOptions(*options.ToGorillaOptions())
}

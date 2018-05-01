package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/gin-contrib/sessions"
	gsessions "github.com/gorilla/sessions"
)

type Store interface {
	sessions.Store
}

// client: memcache client.
// keyPrefix: prefix for the keys we store.
func NewStore(
	client *memcache.Client, keyPrefix string, keyPairs ...[]byte,
) Store {
	return &store{gsm.NewMemcacheStore(client, keyPrefix, keyPairs...)}
}

type store struct {
	*gsm.MemcacheStore
}

func (c *store) Options(options sessions.Options) {
	c.MemcacheStore.Options = &gsessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

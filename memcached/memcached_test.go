package memcached

import (
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
)

const memcachedTestServer = "localhost:11211"

var newStore = func(_ *testing.T) sessions.Store {
	store := NewStore(
		memcache.New(memcachedTestServer), "", []byte("secret"))
	return store
}

func TestMemcached_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestMemcached_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestMemcached_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestMemcached_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestMemcached_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

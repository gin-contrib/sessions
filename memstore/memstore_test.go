package memstore

import (
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
)

var newStore = func(_ *testing.T) sessions.Store {
	store := NewStore([]byte("secret"))
	return store
}

func TestCookie_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestCookie_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestCookie_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestCookie_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestCookie_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

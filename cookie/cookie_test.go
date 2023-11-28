package cookie

import (
	"testing"

	sessions "github.com/geschke/gin-contrib-sessions"
	"github.com/geschke/gin-contrib-sessions/cookie"
	"github.com/geschke/gin-contrib-sessions/tester"
)

var newStore = func(_ *testing.T) sessions.Store {
	store := cookie.NewStore([]byte("secret"))
	return store
}

func TestCookie_SessionGetSet(t *testing.T) {
	t.Logf("in GetSet")
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
	t.Logf("in TestCookieSessionOptions")
	tester.Options(t, newStore)
}

func TestCookie_SessionMany(t *testing.T) {
	tester.Many(t, newStore)
}

package sqlite

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
)

var session_options = sessions.Options{
	Path:     "/",
	Domain:   "localhost",
	MaxAge:   600,
	Secure:   false,
	HttpOnly: false,
	SameSite: http.SameSiteDefaultMode,
}
var newStore = func(t *testing.T) sessions.Store {
	store, err := NewStore("test.db", 10, "sessions", &session_options, []byte("secret-key"))
	if err != nil {
		t.Fatal(err.Error())
	}

	return store
}

func TestSqlite_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
	t.Cleanup(cleanup)
}

func TestSqlite_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
	t.Cleanup(cleanup)
}

func TestSqlite_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
	t.Cleanup(cleanup)
}

func TestSqlite_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
	t.Cleanup(cleanup)
}

func TestSqlite_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
	t.Cleanup(cleanup)
}

func TestSqlite_SessionMany(t *testing.T) {
	tester.Many(t, newStore)
	t.Cleanup(cleanup)
}

func cleanup() {
	os.Remove("test.db")
	os.Remove("test.db-shm")
	os.Remove("test.db-wal")
}

package postgres

import (
	"errors"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
)

const testServer = "postgresql://foo:bar@localhost/test"

var newStore = func(t *testing.T) sessions.Store {
	s, err := NewStore(testServer, []byte("secret"))
	if err != nil {
		t.Error(err)
	}
	if s == nil {
		t.Error(errors.New("store returned nil"))
	}
	return s
}

func TestPostgres_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestPostgres_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestPostgres_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestPostgres_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestPostgres_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

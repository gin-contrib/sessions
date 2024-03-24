package postgres

import (
	"database/sql"
	"testing"

	"github.com/weisskopfjens/sessions"
	"github.com/weisskopfjens/sessions/tester"
)

const postgresTestServer = "postgres://testuser:testpw@localhost:5432/testdb?sslmode=disable"

var newStore = func(_ *testing.T) sessions.Store {
	db, err := sql.Open("postgres", postgresTestServer)
	if err != nil {
		panic(err)
	}

	store, err := NewStore(db, []byte("secret"))
	if err != nil {
		panic(err)
	}

	return store
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

func TestPostgres_SessionMany(t *testing.T) {
	tester.Many(t, newStore)
}

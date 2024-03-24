package mysql

import (
	"database/sql"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
)

const mysqlTestServer = "testuser:testpass_@tcp(localhost:3306)/testdb"

var newStore = func(_ *testing.T) sessions.Store {
	db, err := sql.Open("mysql", mysqlTestServer)
	if err != nil {
		panic(err)
	}

	store, err := NewStore(db, []byte("secret"))
	if err != nil {
		panic(err)
	}

	return store
}

func TestMysql_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestMysql_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestMysql_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestMysql_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestMysql_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

func TestMysql_SessionMany(t *testing.T) {
	tester.Many(t, newStore)
}

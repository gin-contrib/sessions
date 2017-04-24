package sessions

import (
	"testing"
)

const mysqlTestServer = "username:password@tcp(localhost)/database?parseTime=true"

var newMySQLStore = func(_ *testing.T) Store {
	store, err := NewMySQLStore(mysqlTestServer, "sessions", []byte("secret"))
	if err != nil {
		panic(err)
	}
	return store
}

func TestMySQL_SessionGetSet(t *testing.T) {
	sessionGetSet(t, newMySQLStore)
}

func TestMySQL_SessionDeleteKey(t *testing.T) {
	sessionDeleteKey(t, newMySQLStore)
}

func TestMySQL_SessionFlashes(t *testing.T) {
	sessionFlashes(t, newMySQLStore)
}

func TestMySQL_SessionClear(t *testing.T) {
	sessionClear(t, newMySQLStore)
}

func TestMySQL_SessionOptions(t *testing.T) {
	sessionOptions(t, newMySQLStore)
}

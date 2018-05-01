package mongo

import (
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
	mgo "gopkg.in/mgo.v2"
)

const mongoTestServer = "localhost:27017"

var newStore = func(_ *testing.T) sessions.Store {
	session, err := mgo.Dial(mongoTestServer)
	if err != nil {
		panic(err)
	}

	c := session.DB("test").C("sessions")
	return NewStore(c, 3600, true, []byte("secret"))
}

func TestMongo_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestMongo_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestMongo_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestMongo_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestMongo_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

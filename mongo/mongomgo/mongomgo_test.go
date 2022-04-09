package mongomgo

import (
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
	"github.com/globalsign/mgo"
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

func TestMongoMGOMGO_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestMongoMGO_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestMongoMGO_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestMongoMGO_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestMongoMGO_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

func TestMongoMGO_SessionMany(t *testing.T) {
	tester.Many(t, newStore)
}

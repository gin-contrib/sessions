package mongodriver

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
	"go.mongodb.org/mongo-driver/mongo"
)

const mongoTestServer = "mongodb://localhost:27017"

var newStore = func(_ *testing.T) sessions.Store {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoTestServer))
	if err != nil {
		panic(err)
	}

	if err := client.Connect(context.Background()); err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	c := client.Database("test").Collection("sessions")
	return NewStore(c, 3600, true, []byte("secret"))
}

func TestMongoDriver_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestMongoDriver_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestMongoDriver_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestMongoDriver_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestMongoDriver_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

func TestMongoDriver_SessionMany(t *testing.T) {
	tester.Many(t, newStore)
}

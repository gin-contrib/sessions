package officialmongo

import (
	"context"
	"testing"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var newStore = func(_ *testing.T) sessions.Store {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connect := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, connect)
	if err != nil {
		panic(err)
	}
	coll := client.Database("test").Collection("sessions")
	return NewStore(coll, 3600, true, []byte("secret"))
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

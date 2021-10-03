package officialmongo

import (
	mongodbstore "github.com/2-72/gorilla-sessions-mongodb"
	"github.com/gin-contrib/sessions"
	g_sessions "github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection
var keyPairs []byte

// Store 連接 gin session
type Store interface {
	sessions.Store
}

var cfg = mongodbstore.MongoDBStoreConfig{}

// NewStore 連接 mongodb
func NewStore(c *mongo.Collection, maxAge int, ensureTTL bool, k []byte) Store {
	collection = c
	keyPairs = k
	cfg.IndexTTL = ensureTTL
	cfg.SessionOptions = g_sessions.Options{
		MaxAge: maxAge,
	}
	s, err := mongodbstore.NewMongoDBStoreWithConfig(collection, cfg, keyPairs)

	if err != nil {
		panic(err)
	}
	return &store{s}
}

type store struct {
	*mongodbstore.MongoDBStore
}

func (c *store) Options(options sessions.Options) {
	newOptions := options.ToGorillaOptions()
	cfg.SessionOptions = *newOptions
	s, err := mongodbstore.NewMongoDBStoreWithConfig(collection, cfg, keyPairs)
	if err != nil {
		panic(err)
	}
	*c = store{s}
}

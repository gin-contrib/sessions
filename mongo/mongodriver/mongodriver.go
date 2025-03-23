package mongodriver

import (
	"github.com/gin-contrib/sessions"

	"github.com/laziness-coders/mongostore"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ sessions.Store = (*store)(nil)

func NewStore(c *mongo.Collection, maxAge int, ensureTTL bool, keyPairs ...[]byte) sessions.Store {
	return &store{mongostore.NewMongoStore(c, maxAge, ensureTTL, keyPairs...)}
}

type store struct {
	*mongostore.MongoStore
}

func (c *store) Options(options sessions.Options) {
	c.MongoStore.Options = options.ToGorillaOptions()
}

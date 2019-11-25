package rw

import (
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/sessions/tester"
)

const redisTestServer = "localhost:6379"

var newRedisStore = func(_ *testing.T) sessions.Store {
	storeRead, err1 := redis.NewStore(10, "tcp", redisTestServer, "", []byte("secret"))
	if err1 != nil {
		panic(err1)
	}
	storeWrite, err2 := redis.NewStore(10, "tcp", redisTestServer, "", []byte("secret"))
	if err2 != nil {
		panic(err2)
	}

	return NewStore(storeRead, storeWrite)
}

func TestRedis_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newRedisStore)
}

func TestRedis_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newRedisStore)
}

func TestRedis_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newRedisStore)
}

func TestRedis_SessionClear(t *testing.T) {
	tester.Clear(t, newRedisStore)
}

func TestRedis_SessionOptions(t *testing.T) {
	tester.Options(t, newRedisStore)
}

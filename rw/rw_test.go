package rw

import (
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/sessions/tester"
)

const redisTestServer = "localhost:6379"

var newRWStore = func(_ *testing.T) sessions.Store {
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

func TestRW_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newRWStore)
}

func TestRW_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newRWStore)
}

func TestRW_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newRWStore)
}

func TestRW_SessionClear(t *testing.T) {
	tester.Clear(t, newRWStore)
}

func TestRW_SessionOptions(t *testing.T) {
	tester.Options(t, newRWStore)
}

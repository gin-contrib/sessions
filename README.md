# sessions

[![Run CI Lint](https://github.com/gin-contrib/sessions/actions/workflows/lint.yml/badge.svg)](https://github.com/gin-contrib/sessions/actions/workflows/lint.yml)
[![Run Testing](https://github.com/gin-contrib/sessions/actions/workflows/testing.yml/badge.svg)](https://github.com/gin-contrib/sessions/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/gin-contrib/sessions/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-contrib/sessions)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/sessions)](https://goreportcard.com/report/github.com/gin-contrib/sessions)
[![GoDoc](https://godoc.org/github.com/gin-contrib/sessions?status.svg)](https://godoc.org/github.com/gin-contrib/sessions)
[![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin)

Gin middleware for session management with multi-backend support:

- [cookie-based](#cookie-based)
- [Redis](#redis)
- [memcached](#memcached)
- [MongoDB](#mongodb)
- [MongoDB official driver](#mongodb-official-driver)
- [memstore](#memstore)

## Usage

### Start using it

Download and install it:

```bash
go get github.com/gin-contrib/sessions
```

Import it in your code:

```go
import "github.com/gin-contrib/sessions"
```

## Basic Examples

### single session

```go
package main

import (
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/hello", func(c *gin.Context) {
    session := sessions.Default(c)

    if session.Get("hello") != "world" {
      session.Set("hello", "world")
      session.Save()
    }

    c.JSON(200, gin.H{"hello": session.Get("hello")})
  })
  r.Run(":8000")
}
```

### multiple sessions

```go
package main

import (
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  sessionNames := []string{"a", "b"}
  r.Use(sessions.SessionsMany(sessionNames, store))

  r.GET("/hello", func(c *gin.Context) {
    sessionA := sessions.DefaultMany(c, "a")
    sessionB := sessions.DefaultMany(c, "b")

    if sessionA.Get("hello") != "world!" {
      sessionA.Set("hello", "world!")
      sessionA.Save()
    }

    if sessionB.Get("hello") != "world?" {
      sessionB.Set("hello", "world?")
      sessionB.Save()
    }

    c.JSON(200, gin.H{
      "a": sessionA.Get("hello"),
      "b": sessionB.Get("hello"),
    })
  })
  r.Run(":8000")
}
```

## Backend Examples

### cookie-based

```go
package main

import (
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### Redis

```go
package main

import (
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/redis"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### Memcached

#### ASCII Protocol

```go
package main

import (
  "github.com/bradfitz/gomemcache/memcache"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/memcached"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  store := memcached.NewStore(memcache.New("localhost:11211"), "", []byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

#### Binary protocol (with optional SASL authentication)

```go
package main

import (
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/memcached"
  "github.com/gin-gonic/gin"
  "github.com/memcachier/mc"
)

func main() {
  r := gin.Default()
  client := mc.NewMC("localhost:11211", "username", "password")
  store := memcached.NewMemcacheStore(client, "", []byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### MongoDB

```go
package main

import (
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/mongo"
  "github.com/gin-gonic/gin"
  "github.com/globalsign/mgo"
)

func main() {
  r := gin.Default()
  session, err := mgo.Dial("localhost:27017/test")
  if err != nil {
    // handle err
  }

  c := session.DB("").C("sessions")
  store := mongo.NewStore(c, 3600, true, []byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

### mongodb-official-driver

```go
package main

import (
	"context"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions/officialmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connect := options.Client().ApplyURI("localhost:27017")
	client, err := mongo.Connect(ctx, connect)
	if err != nil {
		panic(err)
	}
	coll := client.Database("test").Collection("sessions")
	store := officialmongo.NewStore(coll, 3600, true, []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/hello", func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(200, gin.H{"hello": session.Get("hello")})
	})
	r.Run(":8000")
}

```

### memstore

```go
package main

import (
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/memstore"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  store := memstore.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/incr", func(c *gin.Context) {
    session := sessions.Default(c)
    var count int
    v := session.Get("count")
    if v == nil {
      count = 0
    } else {
      count = v.(int)
      count++
    }
    session.Set("count", count)
    session.Save()
    c.JSON(200, gin.H{"count": count})
  })
  r.Run(":8000")
}
```

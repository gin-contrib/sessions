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

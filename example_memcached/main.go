package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-contrib/sessions"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	r := gin.Default()
	store := sessions.NewMemcacheStore(memcache.New("localhost:11211"), "", []byte("secret"))
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

package main

import (
	"github.com/gin-contrib/sessions"
	redistore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func main() {
	r := gin.Default()
	redisStore := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
	})
	store, _ := redistore.NewStore(redisStore)
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

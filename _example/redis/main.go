package main

import (
	"github.com/gin-contrib/sessions"
	redissessions "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	r := gin.Default()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})
	store := redissessions.NewStore(client, "session_", []byte("secret"))
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

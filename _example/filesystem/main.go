package main

import (
	sessions "github.com/geschke/gin-contrib-sessions"
	"github.com/geschke/gin-contrib-sessions/filesystem"
	"github.com/gin-gonic/gin"
)

func main() {
	sessionPath := "/tmp/"
	r := gin.Default()
	store := filesystem.NewStore(sessionPath, []byte("secret"))
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

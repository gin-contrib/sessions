package main

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/rabex-io/sessions"
	"github.com/rabex-io/sessions/mongo"
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

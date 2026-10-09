package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mysql"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := gin.Default()

	userStr := ""
	passStr := ""
	ipStr := ""
	databaseStr := ""
	portStr := ""
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",userStr,passStr,ipStr,portStr,databaseStr)
	store,err := mysql.NewStore(connStr,"sessionstore", "/", 3600, []byte("kukunet"))
	if err != nil {
		fmt.Println("Error connecting to MySQL: ", err)
	}

	//add session
	r.Use(sessions.Sessions("mySession", &store))
	defer store.Close()
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

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

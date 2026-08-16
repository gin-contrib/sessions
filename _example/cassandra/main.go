package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cassandra"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func main() {
	// Requires a running Cassandra with the keyspace and table:
	//
	//	CREATE KEYSPACE IF NOT EXISTS sessions
	//	  WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
	//	CREATE TABLE IF NOT EXISTS sessions.sessions
	//	  (session_id text PRIMARY KEY, data blob, expires_at timestamp);
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "sessions"
	cluster.Consistency = gocql.Quorum

	store, err := cassandra.NewStore(
		cluster,
		3600,
		[]byte("authentication-key"),
		[]byte("0123456789abcdef"),
	)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	r := gin.Default()
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/login", func(c *gin.Context) {
		var req struct {
			UserID   int    `json:"user_id" form:"user_id" binding:"required"`
			Username string `json:"username" form:"username"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", req.UserID)
		session.Set("username", req.Username)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
			"user_id": req.UserID,
		})
	})

	r.GET("/me", func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": session.Get("username"),
		})
	})

	r.POST("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Options(sessions.Options{Path: "/", MaxAge: -1})
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
	})

	r.Run(":8080")
}

// To test this example, you can run the following commands to set up Cassandra:
// docker run -d --name cassandra -p 9042:9042 cassandra:5    # if not already up
// docker exec -it cassandra cqlsh -e "
// CREATE KEYSPACE IF NOT EXISTS sessions WITH replication = {'class':'SimpleStrategy','replication_factor':1};
// CREATE TABLE IF NOT EXISTS sessions.sessions (session_id text PRIMARY KEY, data blob, expires_at timestamp);"

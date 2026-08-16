//go:build example

package cassandra

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func Example() {
	cluster := gocql.NewCluster(
		"127.0.0.1",
	)

	cluster.Keyspace = "sessions"
	cluster.Consistency = gocql.Quorum

	store, err := NewStore(
		cluster,
		3600,
		[]byte("authentication-key"),
		[]byte("encryption-key"),
	)
	if err != nil {
		panic(err)
	}

	defer store.Close()

	router := gin.Default()

	router.Use(
		sessions.Sessions(
			"mysession",
			store,
		),
	)

	router.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)

		session.Set("user_id", 123)
		session.Set("username", "raza")

		if err := session.Save(); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": err.Error(),
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "login successful",
			},
		)
	})

	router.GET("/me", func(c *gin.Context) {
		session := sessions.Default(c)

		userID := session.Get("user_id")

		if userID == nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "unauthorized",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"user_id":  userID,
				"username": session.Get("username"),
			},
		)
	})

	router.POST("/logout", func(c *gin.Context) {
		session := sessions.Default(c)

		session.Clear()
		session.Options(sessions.Options{
			Path:   "/",
			MaxAge: -1,
		})

		if err := session.Save(); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": err.Error(),
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "logout successful",
			},
		)
	})

	fmt.Println("server listening on :8080")

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

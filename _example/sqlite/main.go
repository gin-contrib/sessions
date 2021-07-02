package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/sqlite"
        "log" 
	"github.com/gin-gonic/gin"
        "net/http"
)

func main() {
	r := gin.Default()
	session_options := sessions.Options{
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   600,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteDefaultMode,
	}
	store, err := sqlite.NewStore("test.db", 10, "sessions", &session_options, []byte("secret-key"))
	r.Use(sessions.Sessions("mysession", store))
	if err != nil {
          log.Fatal("failed to create session store",err)
        }
	r.GET("/hello", func(c *gin.Context) {
		session := sessions.Default(c)
	if session.Get("foo") == nil {
		session.Set("foo","=>")
	}
	session.Set("foo","=" + session.Get("foo").(string))
			session.Save()

		c.String(http.StatusOK, session.Get("foo").(string))
	})
	r.Run(":8000")
}

// +build go1.11

// Test for SameSite
package tester

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func OptionsGo11(t *testing.T, newStore storeFactory) {
	r := gin.Default()
	store := newStore(t)
	store.Options(sessions.Options{
		Domain: "localhost",
	})
	r.Use(sessions.Sessions(sessionName, store))

	r.GET("/path", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		session.Save()
		c.String(200, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/path", nil)
	r.ServeHTTP(res1, req1)

	s := strings.Split(res1.Header().Get("Set-Cookie"), ";")
	if s[1] != " Domain=localhost" {
		t.Error("Error writing domain with options:", s[1])
	}

	if s[2] != " SameSite=Strict" {
		t.Error("Error writing SameSite with options:", s[2])
	}
}

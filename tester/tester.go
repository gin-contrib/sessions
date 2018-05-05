// Package tester is a package to test each packages of session stores, such as
// cookie, redis, memcached, mongo, memstore.  You can use this to test your own session
// stores.
package tester

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type storeFactory func(*testing.T) sessions.Store

const sessionName = "mysession"

const ok = "ok"

func init() {
	gin.SetMode(gin.TestMode)
}

func GetSet(t *testing.T, newStore storeFactory) {
	r := gin.Default()
	r.Use(sessions.Sessions(sessionName, newStore(t)))

	r.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		session.Save()
		c.String(200, ok)
	})

	r.GET("/get", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("key") != ok {
			t.Error("Session writing failed")
		}
		session.Save()
		c.String(200, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/get", nil)
	req2.Header.Set("Cookie", res1.Header().Get("Set-Cookie"))
	r.ServeHTTP(res2, req2)
}

func DeleteKey(t *testing.T, newStore storeFactory) {
	r := gin.Default()
	r.Use(sessions.Sessions(sessionName, newStore(t)))

	r.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		session.Save()
		c.String(200, ok)
	})

	r.GET("/delete", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("key")
		session.Save()
		c.String(200, ok)
	})

	r.GET("/get", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("key") != nil {
			t.Error("Session deleting failed")
		}
		session.Save()
		c.String(200, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/delete", nil)
	req2.Header.Set("Cookie", res1.Header().Get("Set-Cookie"))
	r.ServeHTTP(res2, req2)

	res3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/get", nil)
	req3.Header.Set("Cookie", res2.Header().Get("Set-Cookie"))
	r.ServeHTTP(res3, req3)
}

func Flashes(t *testing.T, newStore storeFactory) {
	r := gin.Default()
	store := newStore(t)
	store.Options(sessions.Options{
		Domain: "localhost",
	})
	r.Use(sessions.Sessions(sessionName, store))

	r.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.AddFlash(ok)
		session.Save()
		c.String(200, ok)
	})

	r.GET("/flash", func(c *gin.Context) {
		session := sessions.Default(c)
		l := len(session.Flashes())
		if l != 1 {
			t.Error("Flashes count does not equal 1. Equals ", l)
		}
		session.Save()
		c.String(200, ok)
	})

	r.GET("/check", func(c *gin.Context) {
		session := sessions.Default(c)
		l := len(session.Flashes())
		if l != 0 {
			t.Error("flashes count is not 0 after reading. Equals ", l)
		}
		session.Save()
		c.String(200, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/flash", nil)
	req2.Header.Set("Cookie", res1.Header().Get("Set-Cookie"))
	r.ServeHTTP(res2, req2)

	res3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/check", nil)
	req3.Header.Set("Cookie", res2.Header().Get("Set-Cookie"))
	r.ServeHTTP(res3, req3)
}

func Clear(t *testing.T, newStore storeFactory) {
	data := map[string]string{
		"key": "val",
		"foo": "bar",
	}
	r := gin.Default()
	store := newStore(t)
	r.Use(sessions.Sessions(sessionName, store))

	r.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		for k, v := range data {
			session.Set(k, v)
		}
		session.Clear()
		session.Save()
		c.String(200, ok)
	})

	r.GET("/check", func(c *gin.Context) {
		session := sessions.Default(c)
		for k, v := range data {
			if session.Get(k) == v {
				t.Fatal("Session clear failed")
			}
		}
		session.Save()
		c.String(200, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/check", nil)
	req2.Header.Set("Cookie", res1.Header().Get("Set-Cookie"))
	r.ServeHTTP(res2, req2)
}

func Options(t *testing.T, newStore storeFactory) {
	r := gin.Default()
	store := newStore(t)
	store.Options(sessions.Options{
		Domain: "localhost",
	})
	r.Use(sessions.Sessions(sessionName, store))

	r.GET("/domain", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		session.Options(sessions.Options{
			Path: "/foo/bar/bat",
		})
		session.Save()
		c.String(200, ok)
	})
	r.GET("/path", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		session.Save()
		c.String(200, ok)
	})
	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/domain", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/path", nil)
	r.ServeHTTP(res2, req2)

	s := strings.Split(res1.Header().Get("Set-Cookie"), ";")
	if s[1] != " Path=/foo/bar/bat" {
		t.Error("Error writing path with options:", s[1])
	}

	s = strings.Split(res2.Header().Get("Set-Cookie"), ";")
	if s[1] != " Domain=localhost" {
		t.Error("Error writing domain with options:", s[1])
	}
}

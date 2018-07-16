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

const key = "key"
const ok = "ok"

func init() {
	gin.SetMode(gin.TestMode)
}

func GetSet(t *testing.T, newStore storeFactory) {
	r := gin.Default()
	r.Use(sessions.Sessions(sessionName, newStore(t)))

	r.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		if err := session.Set(key, ok); err != nil {
			t.Error(err)
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
	})

	r.GET("/get", func(c *gin.Context) {
		session := sessions.Default(c)
		if val, err := session.Get(key); err != nil {
			t.Error(err)
		} else if val != ok {
			t.Error("Session writing failed")
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
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
		if err := session.Set(key, ok); err != nil {
			t.Error(err)
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
	})

	r.GET("/delete", func(c *gin.Context) {
		session := sessions.Default(c)
		if err := session.Delete(key); err != nil {
			t.Error(err)
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
	})

	r.GET("/get", func(c *gin.Context) {
		session := sessions.Default(c)
		if val, err := session.Get(key); err != nil {
			t.Error(err)
		} else if val != nil {
			t.Error("Session deleting failed")
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
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
		if err := session.AddFlash(ok); err != nil {
			t.Error(err)
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
	})

	r.GET("/flash", func(c *gin.Context) {
		session := sessions.Default(c)
		var l int
		if xs, err := session.Flashes(); err != nil {
			t.Error(err)
		} else {
			l = len(xs)
		}
		if l != 1 {
			t.Error("Flashes count does not equal 1. Equals ", l)
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
	})

	r.GET("/check", func(c *gin.Context) {
		session := sessions.Default(c)
		var l int
		if xs, err := session.Flashes(); err != nil {
			t.Error(err)
		} else {
			l = len(xs)
		}
		if l != 0 {
			t.Error("flashes count is not 0 after reading. Equals ", l)
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
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
			if err := session.Set(k, v); err != nil {
				t.Error(err)
			}
		}
		if err := session.Clear(); err != nil {
			t.Error(err)
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
	})

	r.GET("/check", func(c *gin.Context) {
		session := sessions.Default(c)
		for k, v := range data {
			if val, err := session.Get(k); err != nil {
				t.Error(err)
			} else if val == v {
				t.Fatal("Session clear failed")
			}
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
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
		if err := session.Set(key, ok); err != nil {
			t.Error(err)
		}
		if err := session.Options(sessions.Options{
			Path: "/foo/bar/bat",
		}); err != nil {
			t.Error(err)
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
	})
	r.GET("/path", func(c *gin.Context) {
		session := sessions.Default(c)
		if err := session.Set(key, ok); err != nil {
			t.Error(err)
		}
		if err := session.Save(); err != nil {
			t.Error(err)
		}
		c.String(http.StatusOK, ok)
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

//go:build go1.11
// +build go1.11

package tester

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func testOptionSameSitego(t *testing.T, r *gin.Engine) {
	r.GET("/sameSite", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		session.Options(sessions.Options{
			SameSite: http.SameSiteStrictMode,
		})
		_ = session.Save()
		c.String(200, ok)
	})

	res3 := httptest.NewRecorder()
	req3, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/sameSite", nil)
	r.ServeHTTP(res3, req3)

	s := strings.Split(res3.Header().Get("Set-Cookie"), "; ")
	if len(s) < 2 {
		t.Fatal("No SameSite=Strict found in options")
	}
	if s[1] != "SameSite=Strict" {
		t.Fatal("Error writing samesite with options:", s[1])
	}
}

// +build go1.11

package tester

import (
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
	req3, _ := http.NewRequest("GET", "/sameSite", nil)
	r.ServeHTTP(res3, req3)

	s := strings.Split(res3.Header().Get("Set-Cookie"), ";")
	if s[1] != " SameSite=Strict" {
		t.Error("Error writing samesite with options:", s[1])
	}
}

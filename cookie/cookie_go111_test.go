// +build go1.11

package cookie

import (
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
)

func TestCookie_SessionOptionsGo11(t *testing.T) {
	tester.OptionsGo11(t, newStore)
}

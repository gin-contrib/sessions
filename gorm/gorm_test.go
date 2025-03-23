//go:build go1.13
// +build go1.13

package gorm

import (
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/tester"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var newStore = func(_ *testing.T) sessions.Store {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return NewStore(db, true, []byte("secret"))
}

func TestGorm_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestGorm_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestGorm_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestGorm_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestGorm_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

func TestGorm_SessionMany(t *testing.T) {
	tester.Many(t, newStore)
}

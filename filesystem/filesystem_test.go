package filesystem

import (
	"os"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/filesystem"
	"github.com/gin-contrib/sessions/tester"
)

var sessionPath = os.TempDir()

var newStore = func(_ *testing.T) sessions.Store {
	store := filesystem.NewStore(sessionPath, []byte("secret"))
	return store
}

func TestFilesystem_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestFilesystem_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestFilesystem_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestFilesystem_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestFilesystem_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

func TestFilesystem_SessionMany(t *testing.T) {
	tester.Many(t, newStore)
}

package sessions

import (
	"testing"

	gorillasessions "github.com/gorilla/sessions"
)

func newFlashTestSession() *session {
	return &session{
		session: &gorillasessions.Session{
			Values: make(map[interface{}]interface{}),
		},
	}
}

func TestFlashesDoesNotMarkUnchangedSessionWritten(t *testing.T) {
	s := newFlashTestSession()

	if flashes := s.Flashes(); len(flashes) != 0 {
		t.Fatalf("expected no flash messages, got %v", flashes)
	}
	if s.Written() {
		t.Fatal("reading an empty flash queue should not mark the session as modified")
	}
}

func TestFlashesMarksSessionWrittenWhenConsumed(t *testing.T) {
	s := newFlashTestSession()
	// Simulate a flash message loaded from an existing session.
	s.Session().AddFlash("notice")

	flashes := s.Flashes()
	if len(flashes) != 1 || flashes[0] != "notice" {
		t.Fatalf("expected one flash message, got %v", flashes)
	}
	if !s.Written() {
		t.Fatal("consuming a flash message should mark the session as modified")
	}

	if flashes := s.Flashes(); len(flashes) != 0 {
		t.Fatalf("expected flash messages to be consumed, got %v", flashes)
	}
	if !s.Written() {
		t.Fatal("reading an empty queue must not reset an earlier modification")
	}
}

func TestFlashesPreservesPreviousModifications(t *testing.T) {
	s := newFlashTestSession()
	s.Set("user", "alice")

	if flashes := s.Flashes(); len(flashes) != 0 {
		t.Fatalf("expected no flash messages, got %v", flashes)
	}
	if !s.Written() {
		t.Fatal("reading an empty flash queue must preserve unsaved changes")
	}
}

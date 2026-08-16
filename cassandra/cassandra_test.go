package cassandra

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	ginsessions "github.com/gin-contrib/sessions"

	"github.com/gocql/gocql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func newTestCluster() *gocql.ClusterConfig {
	cluster := gocql.NewCluster("127.0.0.1")

	cluster.Keyspace = "sessions"
	cluster.Consistency = gocql.One

	return cluster
}

func TestNewStore(t *testing.T) {
	cluster := newTestCluster()

	store, err := NewStore(
		cluster,
		3600,
		[]byte("authentication-key"),
		[]byte("encryption-key"),
	)

	if err != nil {
		// Cassandra may not be running during unit tests.
		t.Skipf("Cassandra unavailable: %v", err)
	}

	defer store.Close()

	if store == nil {
		t.Fatal("expected store")
	}

	if store.options.MaxAge != 3600 {
		t.Fatalf(
			"expected MaxAge 3600, got %d",
			store.options.MaxAge,
		)
	}

	if store.table != defaultTable {
		t.Fatalf(
			"expected table %q, got %q",
			defaultTable,
			store.table,
		)
	}
}

func TestNewStoreNilCluster(t *testing.T) {
	_, err := NewStore(
		nil,
		3600,
		[]byte("key"),
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewStoreWithoutKeys(t *testing.T) {
	cluster := newTestCluster()

	_, err := NewStore(
		cluster,
		3600,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOptions(t *testing.T) {
	store := &Store{}

	store.Options(ginsessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})

	if store.options == nil {
		t.Fatal("expected options")
	}

	if store.options.MaxAge != 3600 {
		t.Fatalf(
			"expected MaxAge 3600, got %d",
			store.options.MaxAge,
		)
	}
}

func TestOptionsUpdate(t *testing.T) {
	store := &Store{
		options: &sessions.Options{
			Path: "/",
		},
	}

	store.Options(ginsessions.Options{
		Path:     "/api",
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	if store.options.Path != "/api" {
		t.Fatalf(
			"expected /api, got %s",
			store.options.Path,
		)
	}

	if store.options.MaxAge != 7200 {
		t.Fatalf(
			"expected 7200, got %d",
			store.options.MaxAge,
		)
	}

	if !store.options.Secure {
		t.Fatal("expected Secure=true")
	}
}

func TestTable(t *testing.T) {
	store := &Store{
		table: defaultTable,
	}

	store.Table("custom_sessions")

	if store.table != "custom_sessions" {
		t.Fatalf(
			"expected custom_sessions, got %s",
			store.table,
		)
	}
}

func TestTableEmpty(t *testing.T) {
	store := &Store{
		table: defaultTable,
	}

	store.Table("")

	if store.table != defaultTable {
		t.Fatalf(
			"expected %s, got %s",
			defaultTable,
			store.table,
		)
	}
}

func TestCloneOptions(t *testing.T) {
	options := &sessions.Options{
		Path:     "/api",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	cloned := cloneOptions(options)

	if cloned == options {
		t.Fatal("expected a copy")
	}

	if cloned.Path != options.Path {
		t.Fatal("path mismatch")
	}

	if cloned.MaxAge != options.MaxAge {
		t.Fatal("MaxAge mismatch")
	}

	cloned.Path = "/changed"

	if options.Path == "/changed" {
		t.Fatal("changing clone modified original")
	}
}

func TestNewCookie(t *testing.T) {
	options := &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	cookie := newCookie(
		"session",
		"abc123",
		options,
	)

	if cookie.Name != "session" {
		t.Fatalf(
			"expected session, got %s",
			cookie.Name,
		)
	}

	if cookie.Value != "abc123" {
		t.Fatalf(
			"expected abc123, got %s",
			cookie.Value,
		)
	}

	if cookie.MaxAge != 3600 {
		t.Fatalf(
			"expected 3600, got %d",
			cookie.MaxAge,
		)
	}

	if !cookie.HttpOnly {
		t.Fatal("expected HttpOnly")
	}

	if !cookie.Secure {
		t.Fatal("expected Secure")
	}
}

func TestDeleteCookie(t *testing.T) {
	options := &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	cookie := newCookie(
		"session",
		"",
		options,
	)

	if cookie.MaxAge >= 0 {
		t.Fatal("expected negative MaxAge")
	}

	if cookie.Expires.IsZero() {
		t.Fatal("expected expiration time")
	}
}

func TestNewSessionWithoutCookie(t *testing.T) {
	store := &Store{
		codecs: secureCodecs(),
		options: &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
		},
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		nil,
	)

	session, err := store.New(
		request,
		"session",
	)

	if err != nil {
		t.Fatal(err)
	}

	if session == nil {
		t.Fatal("expected session")
	}

	if !session.IsNew {
		t.Fatal("expected new session")
	}
}

// newIntegrationStore provisions the keyspace and table and returns a live
// store. It skips the calling test when Cassandra is unavailable.
func newIntegrationStore(t *testing.T) *Store {
	t.Helper()

	admin := gocql.NewCluster("127.0.0.1")
	admin.Consistency = gocql.One

	adminSession, err := admin.CreateSession()
	if err != nil {
		t.Skipf("Cassandra unavailable: %v", err)
	}

	defer adminSession.Close()

	if err := adminSession.Query(
		`CREATE KEYSPACE IF NOT EXISTS sessions ` +
			`WITH replication = ` +
			`{'class': 'SimpleStrategy', 'replication_factor': 1}`,
	).Exec(); err != nil {
		t.Skipf("cannot create keyspace: %v", err)
	}

	if err := adminSession.Query(
		`CREATE TABLE IF NOT EXISTS sessions.sessions ` +
			`(session_id text PRIMARY KEY, data blob, expires_at timestamp)`,
	).Exec(); err != nil {
		t.Skipf("cannot create table: %v", err)
	}

	store, err := NewStore(
		newTestCluster(),
		3600,
		[]byte("authentication-key"),
		[]byte("0123456789abcdef"),
	)
	if err != nil {
		t.Skipf("Cassandra unavailable: %v", err)
	}

	t.Cleanup(func() {
		_ = store.Clear(context.Background())
		store.Close()
	})

	return store
}

func TestSessionRoundTrip(t *testing.T) {
	store := newIntegrationStore(t)

	recorder := httptest.NewRecorder()

	saveRequest := httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		nil,
	)

	session, err := store.New(saveRequest, "session")
	if err != nil {
		t.Fatal(err)
	}

	session.Values["user_id"] = 123
	session.Values["name"] = "raza"

	if err := store.Save(saveRequest, recorder, session); err != nil {
		t.Fatal(err)
	}

	if session.ID == "" {
		t.Fatal("expected a generated session ID")
	}

	cookies := recorder.Result().Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected a session cookie")
	}

	loadRequest := httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		nil,
	)

	loadRequest.AddCookie(cookies[0])

	loaded, err := store.New(loadRequest, "session")
	if err != nil {
		t.Fatal(err)
	}

	if loaded.IsNew {
		t.Fatal("expected existing session")
	}

	if loaded.ID != session.ID {
		t.Fatalf(
			"expected id %q, got %q",
			session.ID,
			loaded.ID,
		)
	}

	if loaded.Values["user_id"] != 123 {
		t.Fatalf(
			"expected user_id=123, got %v",
			loaded.Values["user_id"],
		)
	}

	if loaded.Values["name"] != "raza" {
		t.Fatalf(
			"expected name=raza, got %v",
			loaded.Values["name"],
		)
	}
}

func TestSessionDelete(t *testing.T) {
	store := newIntegrationStore(t)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		nil,
	)

	session, err := store.New(request, "session")
	if err != nil {
		t.Fatal(err)
	}

	session.Values["user_id"] = 123

	if err := store.Save(request, recorder, session); err != nil {
		t.Fatal(err)
	}

	cookies := recorder.Result().Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected a session cookie")
	}

	// Delete the session.
	deleteRecorder := httptest.NewRecorder()

	session.Options.MaxAge = -1

	if err := store.Save(request, deleteRecorder, session); err != nil {
		t.Fatal(err)
	}

	// The original cookie should no longer resolve to a stored session.
	loadRequest := httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		nil,
	)

	loadRequest.AddCookie(cookies[0])

	loaded, err := store.New(loadRequest, "session")
	if err != nil {
		t.Fatal(err)
	}

	if !loaded.IsNew {
		t.Fatal("expected deleted session to resolve as new")
	}
}

func TestGenerateID(t *testing.T) {
	first, err := generateID()
	if err != nil {
		t.Fatal(err)
	}

	second, err := generateID()
	if err != nil {
		t.Fatal(err)
	}

	if first == "" {
		t.Fatal("expected a non-empty id")
	}

	if first == second {
		t.Fatal("expected unique ids")
	}

	if strings.Contains(first, "=") {
		t.Fatalf("expected padding to be trimmed, got %q", first)
	}
}

func TestNewSessionInvalidCookie(t *testing.T) {
	store := &Store{
		codecs: secureCodecs(),
		options: &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
		},
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		nil,
	)

	request.AddCookie(
		&http.Cookie{
			Name:  "session",
			Value: "invalid-session",
		},
	)

	session, err := store.New(
		request,
		"session",
	)

	if err != nil {
		t.Fatal(err)
	}

	if !session.IsNew {
		t.Fatal("expected invalid cookie to produce new session")
	}
}

func TestCookieExpiration(t *testing.T) {
	options := &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	cookie := newCookie(
		"session",
		"",
		options,
	)

	if cookie.Expires.After(time.Now()) {
		t.Fatal("expected cookie to be expired")
	}
}

func secureCodecs() []securecookie.Codec {
	return securecookie.CodecsFromPairs(
		[]byte("authentication-key"),
		// The block key must be a valid AES key size (16/24/32 bytes).
		[]byte("0123456789abcdef"),
	)
}

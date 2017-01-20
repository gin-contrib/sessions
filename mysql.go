package sessions

import (
	"github.com/srinathgs/mysqlstore"
	"github.com/gorilla/sessions"
)

type MySQLStore interface {
	Store
}

// endpoint: MySQL DSN, "username:password@protocol(address)/dbname?parseTime=true"
// table: Table where sessions are to be saved. Created automatically.
// Key are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// By default, cookie max age is set to 0, and path is set to root "/".
// This can be changed, however:
//
//     session.Options(Options{
//         MaxAge: 3600,
//         Domain: "example.com",
//         Path:   "/foo/bar",
//     })
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewMySQLStore(endpoint string, table string, keys ...[]byte) (MySQLStore, error) {
	store, err := mysqlstore.NewMySQLStore(endpoint, table, "/", 0, keys...)
	if err != nil {
		return nil, err
	}

	return &mysqlStore{store}, nil
}

type mysqlStore struct {
	*mysqlstore.MySQLStore
}

func (c *mysqlStore) Options(options Options) {
	c.MySQLStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

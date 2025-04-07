package redis

import (
	"errors"

	"github.com/gin-contrib/sessions"

	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
)

type Store interface {
	sessions.Store
}

// size: maximum number of idle connections.
// network: tcp or udp
// address: host:port
// password: redis-password
// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewStore(size int, network, address, username, password string, keyPairs ...[]byte) (Store, error) {
	s, err := redistore.NewRediStore(size, network, address, username, password, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

// NewStoreWithDB creates a new Redis-based session store with the specified parameters.
//
// Parameters:
// - size: The maximum number of idle connections in the pool.
// - network: The network type (e.g., "tcp").
// - address: The address of the Redis server (e.g., "localhost:6379").
// - password: The password for the Redis server (if any).
// - DB: The Redis database to be selected after connecting.
// - keyPairs: A variadic list of byte slices used for authentication and encryption.
//
// Returns:
// - Store: The created session store.
// - error: An error if the store could not be created.
func NewStoreWithDB(size int, network, address, username, password, db string, keyPairs ...[]byte) (Store, error) {
	s, err := redistore.NewRediStoreWithDB(size, network, address, username, password, db, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

// NewStoreWithPool creates a new session store using a Redis connection pool.
// It takes a redis.Pool and an optional variadic list of key pairs for
// authentication and encryption of session data.
//
// Parameters:
//   - pool: A redis.Pool object that manages a pool of Redis connections.
//   - keyPairs: Optional variadic list of byte slices used for authentication
//     and encryption of session data.
//
// Returns:
//   - Store: A new session store backed by Redis.
//   - error: An error if the store could not be created.
func NewStoreWithPool(pool *redis.Pool, keyPairs ...[]byte) (Store, error) {
	s, err := redistore.NewRediStoreWithPool(pool, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

type store struct {
	*redistore.RediStore
}

// GetRedisStore retrieves the Redis store from the provided Store interface.
// It returns an error if the provided Store is not of the expected type.
//
// Parameters:
//   - s: The Store interface from which to retrieve the Redis store.
//
// Returns:
//   - err: An error if the provided Store is not of the expected type.
//   - rediStore: The retrieved Redis store, or nil if there was an error.
func GetRedisStore(s Store) (rediStore *redistore.RediStore, err error) {
	realStore, ok := s.(*store)
	if !ok {
		err = errors.New("unable to get the redis store: Store isn't *store")
		return nil, err
	}

	rediStore = realStore.RediStore
	return rediStore, nil
}

// SetKeyPrefix sets a key prefix for the given Redis store.
// It retrieves the Redis store from the provided Store interface and sets the key prefix.
// If there is an error retrieving the Redis store, it returns the error.
//
// Parameters:
//   - s: The Store interface from which the Redis store will be retrieved.
//   - prefix: The key prefix to be set for the Redis store.
//
// Returns:
//   - error: An error if there is an issue retrieving the Redis store, otherwise nil.
func SetKeyPrefix(s Store, prefix string) error {
	rediStore, err := GetRedisStore(s)
	if err != nil {
		return err
	}

	rediStore.SetKeyPrefix(prefix)
	return nil
}

func (c *store) Options(options sessions.Options) {
	c.RediStore.Options = options.ToGorillaOptions()
}

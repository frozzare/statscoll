package cache

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/coreos/bbolt"
)

// Cache represents the cache.
type Cache struct {
	bucket []byte
	db     *bolt.DB
}

// New creates a new cache.
func New() (*Cache, error) {
	db, err := bolt.Open("cache.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return nil, err
	}

	return &Cache{
		bucket: []byte("cache"),
		db:     db,
	}, nil
}

// Close the cache.
func (c *Cache) Close() error {
	return c.db.Close()
}

// Get gets a value from the cache and returns it or a error.
func (c *Cache) Get(key string) (interface{}, error) {
	var v interface{}

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucket)
		buf := b.Get([]byte(key))

		return json.Unmarshal(buf, &v)
	})

	if err != nil {
		return nil, err
	}

	return v, nil
}

// Set sets a value and return a error if any.
func (c *Cache) Set(key string, value interface{}) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(c.bucket)
		if err != nil {
			return err
		}

		buf, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return b.Put([]byte(key), buf)
	})
}

// RemovePrefix removes keys by prefix.
func (c *Cache) RemovePrefix(key string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucket)
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if strings.HasPrefix(string(k), key) {
				if err := b.Delete(k); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
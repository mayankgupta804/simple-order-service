package database

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

type DB struct {
	client *bolt.DB
}

func NewInstance(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	return &DB{client: db}, nil
}

func (db *DB) Put(schema, key, value []byte) error {
	err := db.client.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(schema)
		if err != nil {
			return err
		}
		err = b.Put(key, value)
		return err
	})
	return err
}

func (db *DB) Get(schema, key []byte) []byte {
	var val []byte
	db.client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(schema)
		val = b.Get(key)
		return nil
	})
	return val
}

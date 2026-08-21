package store

import (
	"errors"
	"go.etcd.io/bbolt"
	"path/filepath"
)

var buckets = [][]byte{[]byte("records"), []byte("events"), []byte("workflows"), []byte("attachments")}

type DB struct {
	raw  *bbolt.DB
	path string
}

func Open(path string) (*DB, error) {
	raw, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{NoSync: true})
	if err != nil {
		return nil, err
	}
	d := &DB{raw: raw, path: path}
	if err := d.initialize(); err != nil {
		raw.Close()
		return nil, err
	}
	return d, nil
}

func (d *DB) initialize() error {
	return d.raw.Update(func(tx *bbolt.Tx) error {
		for _, name := range buckets {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}
func (d *DB) Close() error {
	if d == nil || d.raw == nil {
		return nil
	}
	return d.raw.Close()
}
func (d *DB) Path() string { return d.path }
func (d *DB) View(fn func(*bbolt.Tx) error) error {
	if d == nil || d.raw == nil {
		return errors.New("database closed")
	}
	return d.raw.View(fn)
}
func (d *DB) Update(fn func(*bbolt.Tx) error) error {
	if d == nil || d.raw == nil {
		return errors.New("database closed")
	}
	return d.raw.Update(fn)
}
func bucketFor(kind string) []byte {
	switch kind {
	case "record":
		return buckets[0]
	case "event":
		return buckets[1]
	case "workflow":
		return buckets[2]
	default:
		return buckets[3]
	}
}

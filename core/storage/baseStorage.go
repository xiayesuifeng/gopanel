package storage

import (
	"errors"
	"go.etcd.io/bbolt"
	"log"
	"os"
	"path"
)

var baseStorage Storage

type BaseStorage struct {
	db *bbolt.DB
}

func (b *BaseStorage) Get(module, key string) (value []byte, err error) {
	err = b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(module))
		if bucket == nil {
			return errors.New("bucket: " + module + " not exist")
		}

		value = bucket.Get([]byte(key))

		return nil
	})

	return
}

func (b *BaseStorage) Set(module, key string, value []byte) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(module))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), value)
	})
}

func (b *BaseStorage) Has(module, key string) (value bool) {
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(module))
		if bucket == nil {
			return errors.New("bucket: " + module + " not exist")
		}

		value = bucket.Get([]byte(key)) != nil

		return nil
	})

	if err != nil {
		return false
	}

	return
}

func (b *BaseStorage) Delete(module, key string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(module))

		return bucket.Delete([]byte(key))
	})
}

func InitBaseStorage(dataPath string) error {
	if _, err := os.Stat(dataPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dataPath, 0755); err != nil {
				log.Panicln("data dir create failure")
			}
		} else {
			log.Panicln("data dir stat failure, error:", err)
		}
	}

	db, err := bbolt.Open(path.Join(dataPath, "gopanel.db"), 0600, nil)
	if err != nil {
		return err
	}

	baseStorage = &BaseStorage{
		db: db,
	}

	return nil
}

func GetBaseStorage() Storage {
	return baseStorage
}

package db

import (
	"gbc/constant"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger"
)

type DB struct {
	*badger.DB
}

func InitDB() *DB {
	opt := badger.DefaultOptions(constant.DBPath)
	opt.Logger = nil
	os.MkdirAll(filepath.Dir(constant.DBPath), os.ModePerm)
	tempDB, err := badger.Open(opt)
	if err != nil {
		panic(err)
	}
	return &DB{tempDB}
}

func (db *DB) Get(key []byte) []byte {
	var value []byte
	db.View(func(txn *badger.Txn) error {
		ob, err := txn.Get(key)
		if err != nil {
			panic(err)
		}
		ob.Value(func(val []byte) error {
			value = val
			return nil
		})
		return nil
	})
	return value
}

func (db *DB) Set(key []byte, value []byte) {
	//You can't get error!
	db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
}

package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const defaultDatabasePath = "url_database"

// URLDatabase is an interface ressembling leveldb.DB, which is
// used to facilitate dependency injection of mocks in tests.
type URLDatabase interface {
	Close() error
	Delete(key []byte, wo *opt.WriteOptions) error
	Get(key []byte, ro *opt.ReadOptions) (value []byte, err error)
	Has(key []byte, ro *opt.ReadOptions) (ret bool, err error)
	Put(key, value []byte, wo *opt.WriteOptions) error
}

// DatabaseManager is an interface ressembling the top level functions
// of leveldb which is used to faciliate dependency injection of mocks
// in tests.
type DatabaseManager interface {
	OpenFile(path string, o *opt.Options) (db *leveldb.DB, err error)
}

// LevelDBDatabaseManager implements logic specific to the LevelDB database manager.
type LevelDBDatabaseManager struct{}

// OpenFile opens a LevelDB database at a given filepath.
func (m *LevelDBDatabaseManager) OpenFile(path string, o *opt.Options) (*leveldb.DB, error) {
	return leveldb.OpenFile(path, nil)
}

// GetURLDatabase retrieves the URL database.
func GetURLDatabase(databaseManager DatabaseManager) URLDatabase {
	urlDatabase, err := databaseManager.OpenFile(defaultDatabasePath, nil)
	if err != nil {
		errorMessage := fmt.Sprintf("Error: Unable to access URL database: %s", err)
		panic(errorMessage)
	}
	return urlDatabase
}

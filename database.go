package main

import (
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// URLDatabase is an interface resembling leveldb.DB which is
// used to facilitate dependency injection of mocks in tests.
type URLDatabase interface {
	Close() error
	Delete(key []byte, wo *opt.WriteOptions) error
	Get(key []byte, ro *opt.ReadOptions) (value []byte, err error)
	Has(key []byte, ro *opt.ReadOptions) (ret bool, err error)
	Put(key, value []byte, wo *opt.WriteOptions) error
}

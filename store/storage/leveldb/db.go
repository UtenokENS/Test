package leveldb

import (
	"bytes"
	"encoding/binary"

	"cosmossdk.io/store/v2"
	dbm "github.com/cosmos/cosmos-db"
)

var _ store.VersionedDatabase = (*DB)(nil)

type DB struct {
	dbm.DB
}

func (db *DB) storageKey(storeKey string, version uint64, key []byte) []byte {
	var buf bytes.Buffer
	buf.Write([]byte(storeKey))
	binary.BigEndian.PutUint64(buf.Bytes(), version)
	buf.Write(key)
	return buf.Bytes()
}

func (db *DB) Has(storeKey string, version uint64, key []byte) (bool, error) {
	return db.DB.Has(db.storageKey(storeKey, version, key))
}

func (db *DB) Get(storeKey string, version uint64, key []byte) ([]byte, error) {
	return db.DB.Get(db.storageKey(storeKey, version, key))
}

func (db *DB) GetLatestVersion() (uint64, error) {
	// TODO: implement me
	return 0, nil
}

func (db *DB) Set(storeKey string, version uint64, key, value []byte) error {
	return db.DB.Set(db.storageKey(storeKey, version, key), value)
}

func (db *DB) Delete(storeKey string, version uint64, key []byte) error {
	return db.DB.Delete(db.storageKey(storeKey, version, key))
}

func (db *DB) SetLatestVersion(version uint64) error {
	panic("not implemented")
}

func (db *DB) NewIterator(storeKey string, version uint64, start, end []byte) (store.Iterator, error) {
	panic("not implemented")
}

func (db *DB) NewReverseIterator(storeKey string, version uint64, start, end []byte) (store.Iterator, error) {
	panic("not implemented")
}

func (db *DB) NewBatch(version uint64) (store.Batch, error) {
	panic("not implemented")
}

func (db *DB) Close() error {
	panic("not implemented")
}

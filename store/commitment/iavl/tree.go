package iavl

import (
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/store/v2/commitment"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/iavl"
	ics23 "github.com/cosmos/ics23/go"
)

var _ commitment.Database = (*Tree)(nil)

// Tree is a wrapper around iavl.MutableTree.
type Tree struct {
	*iavl.MutableTree
}

func (t *Tree) Set(key, value []byte) (bool, error) {
	return t.MutableTree.Set(key, value)
}

func (t *Tree) Get(key []byte) ([]byte, error) {
	return t.MutableTree.Get(key)
}

func (t *Tree) Remove(key []byte) ([]byte, bool, error) {
	return t.MutableTree.Remove(key)
}

// NewIavlTree creates a new Tree instance.
func NewIavlTree(db dbm.DB, logger log.Logger, cfg *Config) *Tree {
	tree := iavl.NewMutableTree(db, cfg.CacheSize, cfg.SkipFastStorageUpgrade, logger)
	return &Tree{
		MutableTree: tree,
	}
}

// WriteBatch writes a batch of key-value pairs to the database.
func (t *Tree) WriteBatch(batch *commitment.Batch) error {
	for _, kv := range batch.Pairs {
		if kv.Value == nil {
			_, deleted, err := t.Remove(kv.Key)
			if err != nil {
				return err
			}
			if !deleted {
				return fmt.Errorf("failed to delete key %X", kv.Key)
			}
		} else {
			_, err := t.Set(kv.Key, kv.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WorkingHash returns the working hash of the database.
func (t *Tree) WorkingHash() []byte {
	return t.MutableTree.WorkingHash()
}

// LoadVersion loads the state at the given version.
func (t *Tree) LoadVersion(version uint64) error {
	return t.LoadVersionForOverwriting(int64(version))
}

// Commit commits the current state to the database.
func (t *Tree) Commit() ([]byte, error) {
	hash, _, err := t.SaveVersion()
	return hash, err
}

// GetProof returns a proof for the given key and version.
func (t *Tree) GetProof(version uint64, key []byte) (*ics23.CommitmentProof, error) {
	imutableTree, err := t.GetImmutable(int64(version))
	if err != nil {
		return nil, err
	}

	return imutableTree.GetProof(key)
}

// GetLatestVersion returns the latest version of the database.
func (t *Tree) GetLatestVersion() uint64 {
	return uint64(t.Version())
}

// Close closes the iavl tree.
func (t *Tree) Close() error {
	return nil
}

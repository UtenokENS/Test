package commitment

import (
	"io"

	ics23 "github.com/cosmos/ics23/go"
)

// Database is an interface for a commitment layer to support multiple backends.
type Database interface {
	WriteBatch(batch *Batch) error
	WorkingHash() []byte
	GetLatestVersion() uint64
	LoadVersion(targetVersion uint64) error
	Commit() ([]byte, error)
	GetProof(version uint64, key []byte) (*ics23.CommitmentProof, error)
	Set(key, value []byte) (bool, error)
	Get(key []byte) ([]byte, error)
	Remove(key []byte) ([]byte, bool, error)

	io.Closer
}

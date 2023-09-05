package multistore

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"

	"cosmossdk.io/store/v2"
	"cosmossdk.io/store/v2/commitment"
	ics23 "github.com/cosmos/ics23/go"
)

type Store struct {
	ss      store.VersionedDatabase
	sc      map[string]commitment.Database
	version uint64
}

func New(ss store.VersionedDatabase) (*Store, error) {
	latestVersion, err := ss.GetLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	return &Store{
		ss:      ss,
		sc:      make(map[string]commitment.Database),
		version: latestVersion,
	}, nil
}

func (s *Store) Close() (err error) {
	err = errors.Join(err, s.ss.Close())
	for _, sc := range s.sc {
		err = errors.Join(err, sc.Close())
	}

	s.ss = nil
	s.sc = nil
	s.version = 0

	return err
}

func (s *Store) MountSCStore(storeKey string, sc commitment.Database) error {
	if _, ok := s.sc[storeKey]; ok {
		return fmt.Errorf("SC store with key %s already mounted", storeKey)
	}

	s.sc[storeKey] = sc
	return nil
}

func (s *Store) GetLatestVersion() (uint64, error) {
	for _, sc := range s.sc {
		v := sc.GetLatestVersion()
		if v != s.version {
			return 0, fmt.Errorf("latest version mismatch for SC store; %d != %d", v, s.version)
		}
	}

	return s.version, nil
}

func (s *Store) GetProof(storeKey string, version uint64, key []byte) (*ics23.CommitmentProof, error) {
	sc, ok := s.sc[storeKey]
	if !ok {
		return nil, fmt.Errorf("SC store with key %s not mounted", storeKey)
	}

	return sc.GetProof(version, key)
}

func (s *Store) GetSCStore(storeKey string) commitment.Database {
	return s.sc[storeKey]
}

func (s *Store) LoadVersion(version uint64) error {
	panic("not implemented!")
}

func (s *Store) WorkingHash() []byte {
	panic("not implemented!")
}

func (s *Store) Commit() ([]byte, error) {
	// TODO this should look more like the following:
	// https://github.com/cosmos/cosmos-sdk/blob/682a9acd83d9abc52941bbb3b92565f6887967fb/store/rootmulti/store.go#L486-L489
	// -> https://github.com/cosmos/cosmos-sdk/blob/80dd55f79bba8ab675610019a5764470a3e2fef9/store/types/commit_info.go#L30-L38
	//  -> https://github.com/cosmos/cosmos-sdk/blob/d1a337eb7828a676978bbbe502ce4b888351b9b3/store/internal/maps/maps.go#L187
	// TODO: needs tests to ensure storev1 Hash() parity

	var hashes bytes.Buffer
	for _, db := range s.sc {
		var h []byte
		h, err := db.Commit()
		if err != nil {
			return nil, err
		}
		hashes.Write(h)
	}
	scHash := sha256.Sum256(hashes.Bytes())

	return scHash[:], nil
}

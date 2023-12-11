package main

import (
	"fmt"
	"io"
)

type Store struct {
	kvMap map[string]string
	db    io.ReadWriteCloser
}

func NewStore(db io.ReadWriteCloser) *Store {
	kv := make(map[string]string)
	return &Store{kvMap: kv, db: db}
}

// Load loads the in-memory KV map with values from the db.
func (s *Store) Load() error {
	return nil
}

// Save stores the in-memory KV map to the db.
func (s *Store) Save() error {
	return nil
}

// Get retreives the value for the given key from the KV map. An error
// is returned if the key is not found.
func (s *Store) Get(key string) (string, error) {
	val, found := s.kvMap[key]
	if !found {
		return "", fmt.Errorf("key %q not found", key)
	}
	return val, nil
}

// Add stores the KV-pair in the KV map. An error will be returned if the
// key already exists.
func (s *Store) Add(key, val string) error {
	_, found := s.kvMap[key]
	if found {
		return fmt.Errorf("key %q already exists", key)
	}
	s.kvMap[key] = val
	return nil
}

// Update replaces the value of an existing key to a new one. An error is
// returned if the key is not found.
func (s *Store) Update(key, val string) error {
	_, found := s.kvMap[key]
	if !found {
		return fmt.Errorf("key %q not found", key)
	}
	s.kvMap[key] = val
	return nil
}

// Delete removes the given key and its value from the KV map. An error
// is returned if the key is not found.
func (s *Store) Delete(key string) error {
	_, found := s.kvMap[key]
	if !found {
		return fmt.Errorf("key %q not found", key)
	}
	delete(s.kvMap, key)
	return nil
}

package main

import "fmt"

type Store map[string]string

func NewStore() *Store {
	s := make(Store)
	return &s
}

func (s *Store) Get(key string) (string, error) {
	val, found := (*s)[key]
	if !found {
		return "", fmt.Errorf("key %q not found", key)
	}
	return val, nil
}

func (s *Store) Add(key, val string) error {
	_, found := (*s)[key]
	if found {
		return fmt.Errorf("key %q already exists", key)
	}
	(*s)[key] = val
	return nil
}

func (s *Store) Delete(key string) error {
	_, found := (*s)[key]
	if !found {
		return fmt.Errorf("key %q not found", key)
	}
	delete(*s, key)
	return nil
}

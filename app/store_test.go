package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	s := NewStore(nil)

	t.Run("can add a value", func(t *testing.T) {
		err := s.Add("hello", "world")
		assert.NoError(t, err)
	})

	t.Run("error when adding existing key", func(t *testing.T) {
		err := s.Add("exists", "value")
		assert.NoError(t, err)

		err = s.Add("exists", "new value")
		assert.Error(t, err)
	})

	t.Run("can get a value", func(t *testing.T) {
		val, err := s.Get("hello")
		assert.NoError(t, err)
		assert.Equal(t, "world", val)
	})

	t.Run("error when getting non-existent value", func(t *testing.T) {
		_, err := s.Get("world")
		assert.Error(t, err)
	})

	t.Run("can update a value", func(t *testing.T) {
		err := s.Add("mykey", "my value")
		assert.NoError(t, err)

		err = s.Update("mykey", "new value")
		assert.NoError(t, err)

		val, err := s.Get("mykey")
		assert.NoError(t, err)
		assert.Equal(t, "new value", val)
	})

	t.Run("can delete a value", func(t *testing.T) {
		err := s.Add("delete", "me")
		assert.NoError(t, err)

		err = s.Delete("delete")
		assert.NoError(t, err)

		_, err = s.Get("delete")
		assert.Error(t, err)
	})
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	s := NewStore()

	t.Run("can add a value", func(t *testing.T) {
		err := s.Add("hello", "world")
		assert.NoError(t, err)
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

	t.Run("can delete a value", func(t *testing.T) {
		err := s.Add("delete", "me")
		assert.NoError(t, err)

		err = s.Delete("delete")
		assert.NoError(t, err)

		_, err = s.Get("delete")
		assert.Error(t, err)
	})
}

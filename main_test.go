package searchtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchtube(t *testing.T) {
	t.Run("Test search 1", func(t *testing.T) {
		results, err := Search("free software song", 1)
		assert.Nil(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, len(results), 1)
	})
	t.Run("Test search 10", func(t *testing.T) {
		results, err := Search("free software song", 10)
		assert.Nil(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, len(results), 10)
	})
	t.Run("Test search all", func(t *testing.T) {
		results, err := Search("free software song", -1)
		assert.Nil(t, err)
		assert.NotNil(t, results)
		assert.GreaterOrEqual(t, len(results), 1)
	})
}

package searchtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchtube(t *testing.T) {
	validate := func(results []*SearchResult, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, results)
		assert.GreaterOrEqual(t, len(results), 1)
		assert.NotEmpty(t, results[0].Thumbnail)
		assert.NotEmpty(t, results[0].Uploader)
		assert.NotEmpty(t, results[0].Title)
		assert.NotEmpty(t, results[0].Duration)
		assert.NotEmpty(t, results[0].ID)
	}
	t.Run("Test search 1", func(t *testing.T) {
		results, err := Search("free software song", 1)
		validate(results, err)
		assert.Equal(t, len(results), 1)
	})
	t.Run("Test search 10", func(t *testing.T) {
		results, err := Search("free software song", 10)
		validate(results, err)
		assert.Equal(t, len(results), 10)
	})
	t.Run("Test search all", func(t *testing.T) {
		results, err := Search("free software song", -1)
		validate(results, err)
		assert.GreaterOrEqual(t, len(results), 1)
	})
	t.Run("Test search URL", func(t *testing.T) {
		results, err := Search("https://www.youtube.com/watch?v=9sJUDx7iEJw", 1)
		validate(results, err)
		result := results[0]
		assert.False(t, result.Live)
		assert.Equal(t, result.Title, "Richard Stallman Free software Song")
		assert.Equal(t, result.Duration, "2:03")
	})
}

package searchtube

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationParser(t *testing.T) {
	assertDuration := func(expectedDurationStr string, durationStr string) {
		expected, err := time.ParseDuration(expectedDurationStr)
		assert.Nil(t, err)
		result := SearchResult{RawDuration: durationStr}
		duration, err := result.GetDuration()
		assert.Nil(t, err)
		assert.Equal(t, expected, duration)
	}

	assertDuration("2m3s", "2:03")
	assertDuration("1h2m3s", "1:02:03")
}

func TestSearchtube(t *testing.T) {
	validate := func(results []*SearchResult, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, results)
		assert.GreaterOrEqual(t, len(results), 1)
		assert.NotEmpty(t, results[0].Thumbnail)
		assert.NotEmpty(t, results[0].Uploader)
		assert.NotEmpty(t, results[0].Title)
		assert.NotEmpty(t, results[0].RawDuration)
		assert.NotEmpty(t, results[0].ID)
	}

	t.Run("Test search 1", func(t *testing.T) {
		results, err := Search("free software song", 1)
		validate(results, err)
		assert.Equal(t, 1, len(results))
	})

	t.Run("Test search 10", func(t *testing.T) {
		results, err := Search("free software song", 10)
		validate(results, err)
		assert.Len(t, results, 10)
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
		assert.Equal(t, "Richard Stallman Free software Song", result.Title)
		assert.Equal(t, "2:03", result.RawDuration)

		duration, err := result.GetDuration()
		assert.Nil(t, err)

		expectedDuration, err := time.ParseDuration("2m3s")
		assert.Nil(t, err)

		assert.Equal(t, expectedDuration, duration)
	})
}

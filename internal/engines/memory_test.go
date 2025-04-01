package engines_test

import (
	"go-metrics/internal/engines"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_SetAndGet(t *testing.T) {
	// Arrange
	storage := engines.NewMemoryStorage[int, string]()
	setter := engines.NewMemorySetter(storage)
	getter := engines.NewMemoryGetter(storage)

	// Act
	setter.Set(1, "test_value")
	val, exists := getter.Get(1)

	// Assert
	require.True(t, exists)
	assert.Equal(t, "test_value", val)
}

func TestMemoryStorage_GetNonExistentKey(t *testing.T) {
	// Arrange
	storage := engines.NewMemoryStorage[int, string]()
	getter := engines.NewMemoryGetter(storage)

	// Act
	val, exists := getter.Get(1)

	// Assert
	assert.False(t, exists)
	assert.Empty(t, val)
}

func TestMemoryStorage_Range(t *testing.T) {
	// Arrange
	storage := engines.NewMemoryStorage[int, string]()
	setter := engines.NewMemorySetter(storage)
	ranger := engines.NewMemoryRanger(storage)

	// Act
	setter.Set(1, "value1")
	setter.Set(2, "value2")
	setter.Set(3, "value3")

	var keys []int
	ranger.Range(func(k int, v string) bool {
		keys = append(keys, k)
		return true
	})

	// Assert
	assert.ElementsMatch(t, []int{1, 2, 3}, keys)
}

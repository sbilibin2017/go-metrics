package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("Init Info Level", func(t *testing.T) {
		err := Init(LevelInfo)
		assert.NoError(t, err, "Initialization should not return an error")
		Info("Test info message", "key", "value")
		Error("Test info message", "key", "value")
		Sync()

	})

}

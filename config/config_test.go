package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const testConfigFilename = "config_test.json"

func TestNew(t *testing.T) {
	defer removeConfig(t)

	var newApiKey string

	t.Run("Generate a completely new config", func(t *testing.T) {
		cfg, _, err := New(testConfigFilename)
		assert.Nil(t, err)

		newApiKey = cfg.APIKey

		assert.NotEmpty(t, newApiKey)

		assert.Equal(t, 64, len(newApiKey))
	})

	t.Run("Load existing config", func(t *testing.T) {
		cfg, _, err := New(testConfigFilename)
		assert.Nil(t, err)

		assert.NotEmpty(t, cfg.APIKey)

		assert.Equal(t, newApiKey, cfg.APIKey)
	})
}

func removeConfig(t *testing.T) {
	err := os.Remove(testConfigFilename)
	assert.Nil(t, err)
}

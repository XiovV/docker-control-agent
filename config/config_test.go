package config

import (
	"os"
	"testing"
)

const testConfigFilename = "config_test.json"

func TestNew(t *testing.T) {
	defer removeConfig(t)

	var newApiKey string

	t.Run("Generate a completely new config", func(t *testing.T) {
		cfg, err := New(testConfigFilename)
		if err != nil {
			t.Fatal("got error:", err)
		}

		newApiKey = cfg.APIKey

		if newApiKey == "" {
			t.Fatal("api key shouldn't be an empty string")
		}

		if len(newApiKey) != 64 {
			t.Fatalf("wrong api key length; expected %d, got %d", 64, len(newApiKey))
		}
	})

	t.Run("Load existing config", func(t *testing.T) {
		cfg, err := New(testConfigFilename)
		if err != nil {
			t.Fatal("got error:", err)
		}

		if cfg.APIKey == "" {
			t.Fatal("api key shouldn't be an empty string")
		}

		if cfg.APIKey != newApiKey {
			t.Fatalf("loaded invalid key; expected %s, got %s", newApiKey, cfg.APIKey)
		}
	})
}

func removeConfig(t *testing.T) {
	err := os.Remove(testConfigFilename)
	if err != nil {
		t.Fatal("got error:", err)
	}
}
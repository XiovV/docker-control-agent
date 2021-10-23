package config

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	APIKey string `json:"api_key"`
}

func New(filename string) (*Config, error) {
	var apiKeyPlaintext string
	file, err := os.Open(filename)

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(filename)
		if err != nil {
			return nil, err
		}

		randomBytes := make([]byte, 16)

		_, err = rand.Read(randomBytes)
		if err != nil {
			return nil, err
		}

		apiKeyPlaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

		cfg := &Config{APIKey: fmt.Sprintf("%x", sha256.Sum256([]byte(apiKeyPlaintext)))}

		data, err := json.MarshalIndent(cfg, "", "	")
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			return nil, err
		}

		fmt.Println("Your new api key is:", apiKeyPlaintext)
		return cfg, nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var cfg Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		panic(err)
	}

	return &cfg, nil
}

func (c Config) CompareHash(plaintext string) bool {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(plaintext)))

	return c.APIKey == hash
}

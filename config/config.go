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
	ApiKey string `json:"api_key"`
}

func New() Config {
	var apiKeyPlaintext string
	file, err := os.Open("config.json")
	defer file.Close()
	if errors.Is(err, os.ErrNotExist) {
		_, err := os.Create("config.json")
		if err != nil {
			panic(err)
		}

		randomBytes := make([]byte, 16)

		_, err = rand.Read(randomBytes)
		if err != nil {
			panic(err)
		}

		apiKeyPlaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

		cfg := Config{ApiKey: fmt.Sprintf("%x", sha256.Sum256([]byte(apiKeyPlaintext)))}

		data, err := json.MarshalIndent(cfg, "", "	")
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile("config.json", data, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Your new api key is:", apiKeyPlaintext)
		return cfg
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var cfg Config
	json.Unmarshal(bytes, &cfg)

	return cfg
}

func (c Config) CompareHash(plaintext string) bool {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(plaintext)))

	return c.ApiKey == hash
}
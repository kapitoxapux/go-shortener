package config

import (
	"os"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	Base    string `env:"BASE_URL"`
	Path    string `env:"FILE_STORAGE_PATH"`
}

func SetEnvConf(address string, base string, path string) *Config {
	env := new(Config)
	env.Address = address
	env.Base = base
	env.Path = path

	return env
}

func GetStoragePath() string {
	path := os.Getenv("FILE_STORAGE_PATH")

	return path
}

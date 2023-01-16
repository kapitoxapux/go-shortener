package config

import "os"

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	Base    string `env:"BASE_URL"`
}

func SetEnvConf(address string, base string) *Config {
	env := new(Config)

	os.Setenv("SERVER_ADDRESS", address)
	os.Setenv("BASE_URL", base)

	env.Address = os.Getenv("SERVER_ADDRESS")
	env.Base = os.Getenv("BASE_URL")

	return env
}
package config

import (
	"flag"
	"os"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	Base    string `env:"BASE_URL"`
	Path    string `env:"FILE_STORAGE_PATH"`
}

var ServerConfig Config

const defaultServerAdress = "localhost:8080"
const defaultBaseURL = "http://localhost:8080"
const defaultStoragePath = ""

func SetConfig() Config {
	addr := flag.String("a", defaultServerAdress, "SERVER_ADDRESS")
	base := flag.String("b", defaultBaseURL, "BASE_URL")
	path := flag.String("f", defaultStoragePath, "FILE_STORAGE_PATH")
	flag.Parse()

	if os.Getenv("SERVER_ADDRESS") == "" {
		ServerConfig.Address = *addr
	} else {
		ServerConfig.Address = os.Getenv("SERVER_ADDRESS")
	}

	if os.Getenv("BASE_URL") == "" {
		ServerConfig.Base = *base
	} else {
		ServerConfig.Base = os.Getenv("BASE_URL")
	}

	if os.Getenv("FILE_STORAGE_PATH") == "" {
		ServerConfig.Path = *path
	} else {
		ServerConfig.Path = os.Getenv("FILE_STORAGE_PATH")
	}

	return ServerConfig
}

func GetConfigAddress() string {

	return ServerConfig.Address
}

func GetConfigBase() string {

	return ServerConfig.Base
}

func GetConfigPath() string {

	return ServerConfig.Path
}

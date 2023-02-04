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

var Secretkey = []byte("самый секретный секрет")

const defaultServerAdress = "localhost:8080"
const defaultBaseURL = "http://localhost:8080"
const defaultStoragePath = ""

func SetConfig() Config {
	addr := flag.String("a", defaultServerAdress, "SERVER_ADDRESS")
	base := flag.String("b", defaultBaseURL, "BASE_URL")
	path := flag.String("f", defaultStoragePath, "FILE_STORAGE_PATH")
	flag.Parse()

	if address := os.Getenv("SERVER_ADDRESS"); address == "" {
		ServerConfig.Address = *addr
	} else {
		ServerConfig.Address = address
	}

	if base_url := os.Getenv("BASE_URL"); base_url == "" {
		ServerConfig.Base = *base
	} else {
		ServerConfig.Base = base_url
	}

	if storage_path := os.Getenv("FILE_STORAGE_PATH"); storage_path == "" {
		ServerConfig.Path = *path
	} else {
		ServerConfig.Path = storage_path
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

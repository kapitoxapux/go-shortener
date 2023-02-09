package config

import (
	"flag"
	"os"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	Base    string `env:"BASE_URL"`
	Path    string `env:"FILE_STORAGE_PATH"`
	DB      string `env:"DATABASE_DSN"`
}

var ServerConfig Config

var Secretkey = []byte("самый секретный секрет")

const defaultServerAdress = "localhost:8080"
const defaultBaseURL = "http://localhost:8080"
const defaultStoragePath = ""
const defaultStorageDB = "host=localhost port=5432 user=postgres password=tl-wn722n dbname=postgres sslmode=disable"

func SetConfig() Config {
	addr := flag.String("a", defaultServerAdress, "SERVER_ADDRESS")
	base := flag.String("b", defaultBaseURL, "BASE_URL")
	path := flag.String("f", defaultStoragePath, "FILE_STORAGE_PATH")
	db := flag.String("d", defaultStorageDB, "DATABASE_DSN")
	flag.Parse()

	if address := os.Getenv("SERVER_ADDRESS"); address == "" {
		ServerConfig.Address = *addr
	} else {
		ServerConfig.Address = address
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL == "" {
		ServerConfig.Base = *base
	} else {
		ServerConfig.Base = baseURL
	}

	if storagePath := os.Getenv("FILE_STORAGE_PATH"); storagePath == "" {
		ServerConfig.Path = *path
	} else {
		ServerConfig.Path = storagePath
	}

	if storageDB := os.Getenv("DATABASE_DSN"); storageDB == "" {
		ServerConfig.DB = *db
	} else {
		ServerConfig.DB = storageDB
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

func GetStorageDB() string {

	return ServerConfig.DB
}

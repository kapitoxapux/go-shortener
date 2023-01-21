package config

import (
	"os"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	Base    string `env:"BASE_URL"`
	Path    string `env:"FILE_STORAGE_PATH"`
}

// type NetAddress struct {
// 	Host string
// 	Port string
// }

// type BaseUrl struct {
// 	Base string
// }

// type StoragePath struct {
// 	Path string
// }

func SetEnvConf(address string, base string, path string) *Config {
	env := new(Config)

	os.Setenv("SERVER_ADDRESS", address)
	os.Setenv("BASE_URL", base)
	os.Setenv("FILE_STORAGE_PATH", path)

	env.Address = os.Getenv("SERVER_ADDRESS")
	env.Base = os.Getenv("BASE_URL")
	env.Path = os.Getenv("FILE_STORAGE_PATH")

	return env
}

func GetStoragePath() string {
	path := os.Getenv("FILE_STORAGE_PATH")
	return path
}

// func (a NetAddress) String() string {
// 	return a.Host + ":" + a.Port
// }

// func (a *NetAddress) Set(s string) error {
// 	hp := strings.Split(s, ":")

// 	a.Host = hp[0]
// 	a.Port = hp[1]

// 	return nil
// }

// func (b BaseUrl) String() string {
// 	return b.Base
// }

// func (b *BaseUrl) Set(s string) error {
// 	b.Base = s

// 	return nil
// }

// func (f StoragePath) String() string {
// 	return f.Path
// }

// func (f *StoragePath) Set(s string) error {
// 	f.Path = s

// 	return nil
// }

func init() {

	// флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
	// флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL);
	// флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH).

}

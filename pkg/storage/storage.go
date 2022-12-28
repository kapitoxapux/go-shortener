package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

var paths = map[string]*Shorter{}

type Shorter struct {
	Id       string
	LongUrl  string
	ShortUrl string
}

func Shortener(url string) string {
	plainText := []byte(url)
	sha256Hash := sha256.Sum256(plainText)

	return hex.EncodeToString(sha256Hash[:])
}

func SetShort(url string) *Shorter {

	id := Shortener(url)

	shorter := new(Shorter)

	shorter.Id = id
	shorter.ShortUrl = "http://localhost:8080/" + id
	shorter.LongUrl = url

	paths[shorter.Id] = shorter

	return shorter
}

func GetShort(id string) string {
	if paths[id] != nil {
		return paths[id].ShortUrl
	}
	return ""
}

func GetFullUrl(id string) string {
	if paths[id] != nil {
		return paths[id].LongUrl
	}
	return ""
}

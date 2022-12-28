package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

var paths = map[string]*Shorter{}

type Shorter struct {
	ID       string
	LongURL  string
	ShortURL string
}

func Shortener(url string) string {
	plainText := []byte(url)
	sha256Hash := sha256.Sum256(plainText)

	return hex.EncodeToString(sha256Hash[:])
}

func SetShort(url string) *Shorter {

	id := Shortener(url)

	shorter := new(Shorter)

	shorter.ID = id
	shorter.ShortURL = "http://localhost:8080/" + id
	shorter.LongURL = url

	paths[shorter.ID] = shorter

	return shorter
}

func GetShort(id string) string {
	if paths[id] != nil {
		return paths[id].ShortURL
	}
	return ""
}

func GetFullURL(id string) string {
	if paths[id] != nil {
		return paths[id].LongURL
	}
	return ""
}

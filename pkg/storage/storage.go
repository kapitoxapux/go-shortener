package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

var paths = map[string]*Shorter{}

type Shorter struct {
	id       string
	longUrl  string
	shortUrl string
}

func shortener(url string) string {
	plainText := []byte(url)
	sha256Hash := sha256.Sum256(plainText)

	return hex.EncodeToString(sha256Hash[:])
}

func setShort(url string) *Shorter {

	id := shortener(url)

	shorter := new(Shorter)

	shorter.id = id
	shorter.shortUrl = "http://localhost:8080/" + id
	shorter.longUrl = url

	paths[shorter.id] = shorter

	return shorter
}

func getShort(id string) string {
	if paths[id] != nil {
		return paths[id].shortUrl
	}
	return ""
}

func getFullUrl(id string) string {
	if paths[id] != nil {
		return paths[id].longUrl
	}
	return ""
}

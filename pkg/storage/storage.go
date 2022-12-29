package storage

import (
	// "crypto/sha256"
	// "encoding/hex"
	"math/rand"
)

const letterBytes = "_.-~1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var short string = ""

var paths = map[string]*Shorter{}

type Shorter struct {
	ID       string
	LongURL  string
	ShortURL string
}

func Shortener(url string) string {
	// plainText := []byte(url)
	// sha256Hash := sha256.Sum256(plainText)
	// return hex.EncodeToString(sha256Hash[:])

	b := make([]byte, 7)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	if paths[string(b)] != nil {
		return ""
	}
	return string(b)

}

func SetShort(url string) *Shorter {
	short = ""
	for short == "" {
		short = Shortener(url)
	}

	shorter := new(Shorter)

	shorter.ID = short
	shorter.ShortURL = "http://localhost:8080/" + short
	shorter.LongURL = url

	paths[short] = shorter

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

package storage

import (
	"math/big"
	"math/rand"
)

var short string = ""

var paths = map[string]*Shorter{}

type Shorter struct {
	ID       string
	LongURL  string
	ShortURL string
}

func Shortener(url string) string {
	s := make([]byte, 7)
	for i := range s {
		s[i] = url[rand.Intn(len(url))]
	}

	b := new(big.Int).SetBytes(s[2:]).Text(62)
	if paths[b] != nil {
		return ""
	}
	return b
}

func SetShort(link string) *Shorter {
	short = ""
	for short == "" {
		short = Shortener(link)
	}

	shorter := new(Shorter)

	shorter.ID = short
	shorter.ShortURL = "http://localhost:8080/" + short

	shorter.LongURL = link

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

package storage

import (
	"math/big"
	"math/rand"
	"myapp/internal/app/service"
)

type InMemDB struct {
	MemoryDB service.Storage
}

func NewInMemDB() *service.Storage {

	return &service.Storage{
		MemoryDB: storage,
	}
}

var db map[string]*service.Shorter

func Shortener(url string) string {
	s := make([]byte, 7)
	for i := range s {
		s[i] = url[rand.Intn(len(url))]
	}

	b := new(big.Int).SetBytes(s[2:]).Text(62)

	if db[b] != nil {
		return ""
	}

	return b
}

func (s *service.Storage) SetShort(link string) (*service.Shorter, bool) {
	shorter := service.NewShorter()
	duplicate := false

	short := ""
	for short == "" {
		short = Shortener(link)
	}
	shorter.ID = short
	shorter.ShortURL = shorter.BaseURL + short
	shorter.LongURL = link
	shorter.Signer.Sign = service.ShorterSignerSet(short).Sign
	shorter.Signer.SignID = service.ShorterSignerSet(short).SignID

	db[short] = &shorter

	return &shorter, duplicate
}

func (s *service.Storage) GetShort(id string) string {
	shortURL := ""
	if db[id] != nil {

		return db[id].ShortURL
	}

	return shortURL
}

func (s *service.Storage) GetFullURL(id string) string {
	longURL := ""
	if db[id] != nil {

		return db[id].LongURL
	}

	return longURL
}

func (s *service.Storage) GetFullList() map[string]*service.Shorter {
	return db
}

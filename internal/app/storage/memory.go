package storage

import (
	"math/big"
	"math/rand"

	"myapp/internal/app/service"
)

type InMemDB struct {
	db map[string]*service.Shorter
}

func NewInMemDB() *InMemDB {

	return &InMemDB{
		db: make(map[string]*service.Shorter),
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

func (s *InMemDB) SetShort(link string) (*service.Shorter, bool) {
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

	s.db[short] = &shorter

	return &shorter, duplicate
}

func (s *InMemDB) GetShort(id string) string {
	shortURL := ""
	if s.db[id] != nil {

		return s.db[id].ShortURL
	}

	return shortURL
}

func (s *InMemDB) GetFullURL(id string) string {
	longURL := ""
	if s.db[id] != nil {

		return s.db[id].LongURL
	}

	return longURL
}

func (s *InMemDB) GetFullList() map[string]*service.Shorter {
	return s.db
}

// func (s *InMemDB) ConnectionDBCheck() (int, string) {

// 	return 404, ""
// }

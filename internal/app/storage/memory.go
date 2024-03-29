package storage

import (
	"log"
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

func (s *InMemDB) SetShort(link string, data string) (*service.Shorter, bool) {
	shorter := service.NewShorter()
	duplicate := false
	short := ""
	for short == "" {
		short = Shortener(link)
	}
	sign := service.ShorterSignerSet(data)
	shorter.ID = short
	shorter.ShortURL = shorter.BaseURL + short
	shorter.LongURL = link
	shorter.Signer.Sign = sign.Sign
	shorter.Signer.ID = sign.ID
	s.db[short] = &shorter

	return &shorter, duplicate
}

func (s *InMemDB) GetShort(id string) string {
	shortURL := ""
	if s.db[id] != nil && s.db[id].Removed == uint8(0) {

		return s.db[id].ShortURL
	}
	if s.db[id] != nil && s.db[id].Removed == uint8(1) {

		return "402"
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

func (s *InMemDB) GetShorter(id string) *service.Shorter {
	shorter := service.NewShorter()
	if s.db[id] == nil {
		log.Fatal("Произошла ошибка получения структуры")
	} else {
		shorter.ID = s.db[id].ID
		shorter.ShortURL = s.db[id].ShortURL
		shorter.LongURL = s.db[id].LongURL
		shorter.Signer.Sign = s.db[id].Signer.Sign
		shorter.Signer.ID = s.db[id].Signer.ID
		shorter.Removed = s.db[id].Removed
	}
	return &shorter
}

func (s *InMemDB) RemoveShorts(list []string) {
	for _, id := range list {
		s.db[id].Removed = uint8(1)
	}
}

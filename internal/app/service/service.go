package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"myapp/internal/app/config"
)

type Signer struct {
	ID   uint32 `json:"signID"`
	Sign []byte `json:"sign"`
}

type Shorter struct {
	ID       string
	LongURL  string `json:"longURL"`
	ShortURL string `json:"shortURL"`
	BaseURL  string `json:"baseURL"`
	Removed  uint8  `json:"isDeleted"`
	Signer
}

type Storage interface {
	SetShort(link string, cookie string) (*Shorter, bool)
	GetShort(id string) string
	GetShorter(id string) *Shorter
	GetFullURL(id string) string
	GetFullList() map[string]*Shorter
	RemoveShorts(list []string)
}

type Service struct {
	Storage Storage
}

func NewService(storage Storage) *Service {

	return &Service{
		Storage: storage,
	}
}

func NewShorter() Shorter {
	shorter := Shorter{}
	shorter.ID = ""
	shorter.LongURL = ""
	shorter.ShortURL = ""
	shorter.BaseURL = config.GetConfigBase() + "/"
	shorter.Removed = 0
	shorter.Signer.ID = 0
	shorter.Signer.Sign = nil

	return shorter
}

func ShorterSignerSet(data string) Signer {
	resource, _ := hex.DecodeString(data)
	id := binary.BigEndian.Uint32(resource[:4])
	h := hmac.New(sha256.New, config.Secretkey)
	h.Write(resource)
	sign := h.Sum(nil)

	return Signer{id, sign}
}

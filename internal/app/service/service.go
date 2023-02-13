package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"myapp/internal/app/config"
)

type Signer struct {
	SignID uint32 `json:"signID"`
	Sign   []byte `json:"sign"`
}

type Shorter struct {
	ID       string
	LongURL  string `json:"longURL"`
	ShortURL string `json:"shortURL"`
	BaseURL  string `json:"baseURL"`
	Signer
}

type Storage interface {
	SetShort(link string) (*Shorter, bool)
	GetShort(id string) string
	GetFullURL(id string) string
	GetFullList() map[string]*Shorter
	ConnectionDBCheck() (int, string)
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
	shorter.Signer.SignID = 0
	shorter.Signer.Sign = nil

	return shorter
}

func ShorterSignerSet(short string) Signer {
	data, _ := hex.DecodeString(short)
	h := hmac.New(sha256.New, config.Secretkey)
	h.Write(data)
	sign := h.Sum(nil)
	id := binary.BigEndian.Uint32(sign[:4])

	return Signer{id, sign}
}

func ConnectionDBCheck() (int, string) {
	db, err := sql.Open("pgx", config.GetStorageDB())
	if err != nil {

		return 500, err.Error()
	}

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	if err != nil {

		return 500, err.Error()
	}

	return 200, ""
}

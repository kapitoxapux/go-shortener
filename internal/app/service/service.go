package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"sync"

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

type Channel struct {
	InputChannel chan *Shorter
}

func NewService(storage Storage) *Service {

	return &Service{
		Storage: storage,
	}
}

func NewListener(inputCh chan *Shorter) *Channel {

	return &Channel{
		InputChannel: inputCh,
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

func FanOut(inputCh chan *Shorter, n int) []chan *Shorter {
	chs := make([]chan *Shorter, 0, n)
	for i := 0; i < n; i++ {
		ch := make(chan *Shorter)
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan *Shorter) {
			for _, ch := range chs {
				close(ch)
			}
		}(chs)

		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}

			list, ok := <-inputCh
			if !ok {
				return
			}
			ch := chs[i]
			ch <- list
		}
	}()

	return chs
}

func FanIn(inputChs ...chan *Shorter) chan string {
	outCh := make(chan string)
	go func() {
		wg := &sync.WaitGroup{}
		for _, inputCh := range inputChs {
			wg.Add(1)
			go func(inputCh chan *Shorter) {
				defer wg.Done()
				for shorter := range inputCh {
					outCh <- shorter.ID
				}
			}(inputCh)
		}
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func NewWorker(input, out chan *Shorter) {
	go func() {
		for shorter := range input {
			out <- shorter
		}
		close(out)
	}()
}

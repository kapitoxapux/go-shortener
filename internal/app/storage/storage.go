package storage

import (
	"bufio"
	"encoding/json"
	"math/big"
	"math/rand"
	"myapp/internal/app/config"
	"os"
)

var short string = ""

var paths = map[string]*Shorter{}

type saver struct {
	file   *os.File
	writer *bufio.Writer
}

type loader struct {
	file    *os.File
	scanner *bufio.Scanner
}

type Shorter struct {
	ID       string
	LongURL  string `json:"longURL"`
	ShortURL string `json:"shortURL"`
	BaseURL  string `json:"baseURL"`
}

func NewShorter() Shorter {
	shorter := Shorter{}
	shorter.ID = ""
	shorter.LongURL = ""
	shorter.ShortURL = ""
	shorter.BaseURL = config.GetConfigBase() + "/"

	return shorter
}

func NewSaver(filename string) (*saver, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &saver{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *saver) WriteShort(shorter *Shorter) error {
	json, err := json.Marshal(&shorter)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(json); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

func (p *saver) Close() error {

	return p.file.Close()
}

func NewReader(filename string) (*loader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &loader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *loader) Close() error {

	return c.file.Close()
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
	pathStorage := config.GetConfigPath()
	shorter := NewShorter()
	if pathStorage == "" {
		short = ""
		for short == "" {
			short = Shortener(link)
		}
		shorter.ID = short
		shorter.ShortURL = shorter.BaseURL + short
		shorter.LongURL = link
		paths[short] = &shorter
	} else {
		reader, _ := NewReader(pathStorage)
		defer reader.Close()

		for reader.scanner.Scan() {
			data := reader.scanner.Bytes()
			_ = json.Unmarshal(data, &shorter)
			if link == shorter.LongURL {
				return &shorter
			}

		}

		saver, _ := NewSaver(pathStorage)
		defer saver.Close()

		short = ""
		for short == "" {
			short = Shortener(link)
		}
		shorter.ID = short
		shorter.ShortURL = shorter.BaseURL + short
		shorter.LongURL = link
		_ = saver.WriteShort(&shorter)
	}

	return &shorter
}

func GetShort(id string) string {
	shortURL := ""
	pathStorage := config.GetConfigPath()
	if pathStorage == "" {
		if paths[id] != nil {

			return paths[id].ShortURL
		}

	} else {
		reader, _ := NewReader(pathStorage)
		defer reader.Close()

		shorter := NewShorter()
		for reader.scanner.Scan() {
			data := reader.scanner.Bytes()

			_ = json.Unmarshal(data, &shorter)
			if id == shorter.ID {
				return shorter.ShortURL
			}

		}

	}

	return shortURL
}

func GetFullURL(id string) string {
	longURL := ""
	pathStorage := config.GetConfigPath()
	if pathStorage == "" {
		if paths[id] != nil {
			return paths[id].LongURL
		}

	} else {
		reader, _ := NewReader(pathStorage)
		defer reader.Close()

		shorter := NewShorter()
		for reader.scanner.Scan() {
			data := reader.scanner.Bytes()

			_ = json.Unmarshal(data, &shorter)
			if id == shorter.ID {
				return shorter.LongURL
			}

		}

	}

	return longURL
}

package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"myapp/internal/app/config"
	"myapp/internal/app/service"
)

type FileDB struct {
	pathStorage string
}

func NewFileDB() *FileDB {
	pathStorage := config.GetConfigPath()

	return &FileDB{
		pathStorage: pathStorage,
	}
}

type saver struct {
	file   *os.File
	writer *bufio.Writer
}

type loader struct {
	file    *os.File
	scanner *bufio.Scanner
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

func (p *saver) WriteShort(shorter *service.Shorter) error {
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

func (s *FileDB) SetShort(link string) (*service.Shorter, bool) {
	shorter := service.NewShorter()
	duplicate := false
	reader, _ := NewReader(s.pathStorage)
	defer reader.Close()
	for reader.scanner.Scan() {
		data := reader.scanner.Bytes()
		_ = json.Unmarshal(data, &shorter)
		if link == shorter.LongURL {
			return &shorter, duplicate
		}

	}
	saver, _ := NewSaver(s.pathStorage)
	defer saver.Close()
	short := ""
	for short == "" {
		short = Shortener(link)
	}
	shorter.ID = short
	shorter.ShortURL = shorter.BaseURL + short
	shorter.LongURL = link
	shorter.Signer.Sign = service.ShorterSignerSet(short).Sign
	shorter.Signer.SignID = service.ShorterSignerSet(short).SignID
	_ = saver.WriteShort(&shorter)

	return &shorter, duplicate
}

func (s *FileDB) GetShort(id string) string {
	shortURL := ""
	reader, _ := NewReader(s.pathStorage)
	defer reader.Close()
	shorter := service.NewShorter()
	for reader.scanner.Scan() {
		data := reader.scanner.Bytes()
		_ = json.Unmarshal(data, &shorter)
		if id == shorter.ID {
			return shorter.ShortURL
		}

	}

	return shortURL
}

func (s *FileDB) GetFullURL(id string) string {
	longURL := ""
	reader, _ := NewReader(s.pathStorage)
	defer reader.Close()
	shorter := service.NewShorter()
	for reader.scanner.Scan() {
		data := reader.scanner.Bytes()
		_ = json.Unmarshal(data, &shorter)
		if id == shorter.ID {
			return shorter.LongURL
		}

	}

	return longURL
}

func (s *FileDB) GetFullList() map[string]*service.Shorter {
	reader, _ := NewReader(s.pathStorage)
	defer reader.Close()
	paths := map[string]*service.Shorter{}
	for reader.scanner.Scan() {
		data := reader.scanner.Bytes()
		shorter := service.NewShorter()
		_ = json.Unmarshal(data, &shorter)
		paths[shorter.ID] = &shorter
	}

	return paths
}
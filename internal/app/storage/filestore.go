package storage

import (
	"bufio"
	"encoding/json"
	"myapp/internal/app/config"
	"myapp/internal/app/service"
	"os"
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

var pathStorage = config.GetConfigPath()

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

func (s *FileDB) GetShort(id string) string {
	shortURL := ""
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

	return shortURL
}

func (s *FileDB) GetFullURL(id string) string {
	longURL := ""
	reader, _ := NewReader(pathStorage)
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
	reader, _ := NewReader(pathStorage)
	defer reader.Close()

	for reader.scanner.Scan() {
		data := reader.scanner.Bytes()

		shorter := service.NewShorter()
		_ = json.Unmarshal(data, &shorter)
		db[shorter.ID] = &shorter
	}

	return db
}

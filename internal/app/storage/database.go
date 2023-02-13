package storage

import (
	"log"
	"myapp/internal/app/config"
	"myapp/internal/app/repository"
	"myapp/internal/app/service"
)

type DB struct {
	repo repository.Repository
}

func NewDB() *service.Storage {

	repo := repository.NewRepository(config.GetStorageDB())

	return &service.Storage{
		repo: repo,
	}
}

func (s service.Storage) SetShort(link string) (*service.Shorter, bool) {
	shorter := service.NewShorter()
	duplicate := false

	shorter = s.repo.ShowShortenerByLong(link)

	return &shorter, duplicate
}

func (s *service.Storage) GetShort(id string) string {
	shortURL := ""
	if result, err := s.repo.ShowShortener(id); err != nil {
		log.Fatal("Короткая ссылка не найдена, произошла ошибка: %w", err)
	} else {
		shortURL = result.ShortURL
	}

	return shortURL
}

func (s *service.Storage) GetFullURL(id string) string {
	longURL := ""
	if result, err := s.repo.ShowShortener(id); err != nil {
		log.Fatal("Полная ссылка не найдена, произошла ошибка: %w", err)
	} else {
		longURL = result.LongURL
	}

	return longURL
}

func (s *service.Storage) GetFullList() map[string]*service.Shorter {
	if results, err := s.repo.ShowShorteners(); err != nil {
		log.Fatal("Произошла ошибка получения списка: %w", err)
	} else {
		for _, model := range results {
			shorter := service.NewShorter()
			shorter.ID = model.ID
			shorter.ShortURL = model.ShortURL
			shorter.LongURL = model.LongURL
			shorter.Sign = model.Sign
			shorter.SignID = model.SignID

			db[model.ID] = &shorter
		}
	}

	return db
}

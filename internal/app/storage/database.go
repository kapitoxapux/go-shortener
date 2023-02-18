package storage

import (
	"fmt"
	"log"
	"time"

	"myapp/internal/app/config"
	"myapp/internal/app/models"
	"myapp/internal/app/repository"
	"myapp/internal/app/service"
)

type DB struct {
	repo repository.Repository
}

func NewDB() *DB {
	repo := repository.NewRepository(config.GetStorageDB())

	return &DB{
		repo: repo,
	}
}

func (db *DB) SetShort(link string) (*service.Shorter, bool) {
	shorter := service.NewShorter()
	duplicate := false
	model, state := db.repo.ShowShortenerByLong(link)
	if state == "Model found" {
		duplicate = true
		shorter.ID = model.ID
		shorter.ShortURL = model.ShortURL
		shorter.LongURL = model.LongURL
		shorter.Signer.Sign = model.Sign
		shorter.Signer.SignID = model.SignID
	} else {
		s := &models.Link{}
		short := Shortener(link)
		s.ID = short
		s.ShortURL = shorter.BaseURL + short
		s.LongURL = link
		s.Sign = service.ShorterSignerSet(short).Sign
		s.SignID = service.ShorterSignerSet(short).SignID
		s.CreatedAt = time.Now()
		m, err := db.repo.CreateShortener(s)
		if err != nil {
			log.Fatal("Model saving repository failed %w", err.Error())
		}
		shorter.ID = m.ID
		shorter.ShortURL = m.ShortURL
		shorter.LongURL = m.LongURL
		shorter.Signer.Sign = m.Sign
		shorter.Signer.SignID = m.SignID
	}

	return &shorter, duplicate
}

func (db *DB) GetShort(id string) string {
	if result, err := db.repo.ShowShortener(id); err != nil {
		return "402"
	} else {
		return result.ShortURL
	}
}

func (db *DB) GetFullURL(id string) string {
	longURL := ""
	if result, err := db.repo.ShowShortener(id); err != nil {
		fmt.Println("Полная ссылка не найдена, произошла ошибка: %w", err)
	} else {
		longURL = result.LongURL
	}

	return longURL
}

func (db *DB) GetFullList() map[string]*service.Shorter {
	paths := map[string]*service.Shorter{}
	if results, err := db.repo.ShowShorteners(); err != nil {
		log.Fatal("Произошла ошибка получения списка: %w", err)
	} else {
		for _, model := range results {
			shorter := service.NewShorter()
			shorter.ID = model.ID
			shorter.ShortURL = model.ShortURL
			shorter.LongURL = model.LongURL
			shorter.Sign = model.Sign
			shorter.SignID = model.SignID
			paths[model.ID] = &shorter
		}
	}

	return paths
}

func (db *DB) GetShorter(id string) *service.Shorter {
	shorter := service.NewShorter()
	if model, err := db.repo.ShowShortenerByID(id); err != nil {
		log.Fatal("Произошла ошибка получения модели: %w", err)
	} else {
		shorter.ID = model.ID
		shorter.ShortURL = model.ShortURL
		shorter.LongURL = model.LongURL
		shorter.Sign = model.Sign
		shorter.SignID = model.SignID
	}

	return &shorter
}

func (db *DB) RemoveShorts(list []string) {
	if err := db.repo.RemoveShorts(list); err != nil {
		log.Fatal("Произошла ошибка удаления: %w", err)
	}

}

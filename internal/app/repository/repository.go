package repository

import (
	"log"
	"myapp/internal/app/models"

	// "myapp/internal/app/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	CreateShortener(m *models.Shortener) (*models.Shortener, error)
	ShowShortener(id string) (*models.Shortener, error)
	ShowShorteners() ([]models.Shortener, error)
	ShowShortenerByLong(link string) (*models.Shortener, string)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) CreateShortener(m *models.Shortener) (*models.Shortener, error) {
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) ShowShortener(id string) (*models.Shortener, error) {

	model := &models.Shortener{}
	if err := r.db.First(model, "id = ?", []byte(id)).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *repository) ShowShorteners() ([]models.Shortener, error) {
	models := []models.Shortener{}
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

func (r *repository) ShowShortenerBySign(m *models.Shortener) (*models.Shortener, error) {
	if err := r.db.Select(m).Where(m.Sign).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) ShowShortenerByLong(link string) (*models.Shortener, string) {
	model := &models.Shortener{}
	if err := r.db.First(model, "long_url = ?", []byte(link)).Error; err != nil {
		return nil, err.Error()
	}
	return model, "Model found"
}

func NewRepository(dns string) Repository {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Gorm repository failed %w", err.Error())
	}

	if exist := db.Migrator().HasTable(&models.Shortener{}); !exist {
		db.Migrator().CreateTable(&models.Shortener{})
	}

	return &repository{db}
}

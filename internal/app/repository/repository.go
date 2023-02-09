package repository

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"myapp/internal/app/models"
)

type Repository interface {
	CreateShortener(m *models.Links) (*models.Links, error)
	ShowShortener(id string) (*models.Links, error)
	ShowShorteners() ([]models.Links, error)
	ShowShortenerByLong(link string) (*models.Links, string)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) CreateShortener(m *models.Links) (*models.Links, error) {
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) ShowShortener(id string) (*models.Links, error) {

	model := &models.Links{}
	if err := r.db.First(model, "id = ?", []byte(id)).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *repository) ShowShorteners() ([]models.Links, error) {
	models := []models.Links{}
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

func (r *repository) ShowShortenerBySign(m *models.Links) (*models.Links, error) {
	if err := r.db.Select(m).Where(m.Sign).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) ShowShortenerByLong(link string) (*models.Links, string) {
	model := &models.Links{}
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

	if exist := db.Migrator().HasTable(&models.Links{}); !exist {
		db.Migrator().CreateTable(&models.Links{})
	}

	return &repository{db}
}

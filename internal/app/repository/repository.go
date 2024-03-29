package repository

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"myapp/internal/app/models"
)

type Repository interface {
	CreateShortener(m *models.Link) (*models.Link, error)
	ShowShortener(id string) (*models.Link, error)
	ShowShorteners() ([]models.Link, error)
	ShowShortenerByLong(link string) *models.Link
	ShowShortenerByID(id string) (*models.Link, error)
	RemoveShorts(list []string) error
}

type repository struct {
	db *gorm.DB
}

func (r *repository) CreateShortener(m *models.Link) (*models.Link, error) {
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) ShowShortener(id string) (*models.Link, error) {
	model := &models.Link{}
	if err := r.db.First(model, "id = ?", []byte(id)).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (r *repository) ShowShorteners() ([]models.Link, error) {
	models := []models.Link{}
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

func (r *repository) ShowShortenerBySign(m *models.Link) (*models.Link, error) {
	if err := r.db.Select(m).Where(m.Sign).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) ShowShortenerByLong(link string) *models.Link {
	model := &models.Link{}
	if err := r.db.Limit(1).Find(model, "long_url = ?", link).Error; err != nil {
		return nil
	}
	return model
}

func (r *repository) ShowShortenerByID(id string) (*models.Link, error) {
	model := &models.Link{}
	if err := r.db.Where("id=?", id).First(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *repository) RemoveShorts(list []string) error {
	model := &models.Link{}
	if err := r.db.Model(model).Where("id IN ?", list).Updates(models.Link{IsDeleted: uint8(1)}).Error; err != nil {
		return err
	}
	return nil
}

func NewRepository(dns string) Repository {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Gorm repository failed %w", err.Error())
	}
	if exist := db.Migrator().HasTable(&models.Link{}); !exist {
		db.Migrator().CreateTable(&models.Link{})
	}
	if d := db.Migrator().HasColumn(&models.Link{}, "IsDeleted"); !d {
		db.Migrator().AddColumn(&models.Link{}, "IsDeleted")
	}

	return &repository{db}
}

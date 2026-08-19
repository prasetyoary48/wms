package repository

import (
	"github.com/prasetyoary48/wms/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(p *domain.Product) error
	FindByID(id uint) (*domain.Product, error)
	FindBySKU(sku string) (*domain.Product, error)
	FindAll() ([]domain.Product, error)
	Update(p *domain.Product) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(p *domain.Product) error {
	return r.db.Create(p).Error
}

func (r *productRepository) FindByID(id uint) (*domain.Product, error) {
	var p domain.Product
	err := r.db.Preload("Category").First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) FindBySKU(sku string) (*domain.Product, error) {
	var p domain.Product
	err := r.db.Where("sku = ?", sku).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) FindAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Preload("Category").Order("id asc").Find(&products).Error
	return products, err
}

func (r *productRepository) Update(p *domain.Product) error {
	return r.db.Save(p).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}
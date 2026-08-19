package repository

import (
	"github.com/prasetyoary48/wms/internal/domain"
	"gorm.io/gorm"
)

type StockMovementRepository interface {
	Create(tx *gorm.DB, m *domain.StockMovement) error
	FindByID(id uint) (*domain.StockMovement, error)
	FindAll() ([]domain.StockMovement, error)
	FindByStatus(status string) ([]domain.StockMovement, error)
	FindByUser(userID uint) ([]domain.StockMovement, error)
	Update(m *domain.StockMovement) error
}

type stockMovementRepository struct{ db *gorm.DB }

func NewStockMovementRepository(db *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) Create(tx *gorm.DB, m *domain.StockMovement) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(m).Error
}

func (r *stockMovementRepository) FindByID(id uint) (*domain.StockMovement, error) {
	var m domain.StockMovement
	err := r.db.Preload("Product").Preload("FromLocation").Preload("ToLocation").
		Preload("Creator").Preload("Approver").First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *stockMovementRepository) FindAll() ([]domain.StockMovement, error) {
	var list []domain.StockMovement
	err := r.db.Preload("Product").Preload("FromLocation").Preload("ToLocation").
		Preload("Creator").Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *stockMovementRepository) FindByStatus(status string) ([]domain.StockMovement, error) {
	var list []domain.StockMovement
	err := r.db.Preload("Product").Preload("FromLocation").Preload("ToLocation").
		Preload("Creator").Where("status = ?", status).Order("created_at asc").Find(&list).Error
	return list, err
}

func (r *stockMovementRepository) FindByUser(userID uint) ([]domain.StockMovement, error) {
	var list []domain.StockMovement
	err := r.db.Preload("Product").Where("created_by = ?", userID).
		Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *stockMovementRepository) Update(m *domain.StockMovement) error {
	return r.db.Save(m).Error
}
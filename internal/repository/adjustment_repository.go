package repository

import (
	"github.com/prasetyoary48/wms/internal/domain"
	"gorm.io/gorm"
)

type AdjustmentRepository interface {
	Create(a *domain.AdjustmentRequest) error
	FindByID(id uint) (*domain.AdjustmentRequest, error)
	FindByStatus(status string) ([]domain.AdjustmentRequest, error)
	FindByUser(userID uint) ([]domain.AdjustmentRequest, error)
	FindAll() ([]domain.AdjustmentRequest, error)
	Update(a *domain.AdjustmentRequest) error
}

type adjustmentRepository struct{ db *gorm.DB }

func NewAdjustmentRepository(db *gorm.DB) AdjustmentRepository {
	return &adjustmentRepository{db: db}
}

func (r *adjustmentRepository) Create(a *domain.AdjustmentRequest) error {
	return r.db.Create(a).Error
}

func (r *adjustmentRepository) FindByID(id uint) (*domain.AdjustmentRequest, error) {
	var a domain.AdjustmentRequest
	err := r.db.Preload("Product").Preload("Location").Preload("Requester").
		Preload("Approver").First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *adjustmentRepository) FindByStatus(status string) ([]domain.AdjustmentRequest, error) {
	var list []domain.AdjustmentRequest
	err := r.db.Preload("Product").Preload("Location").Preload("Requester").
		Where("status = ?", status).Order("created_at asc").Find(&list).Error
	return list, err
}

func (r *adjustmentRepository) FindByUser(userID uint) ([]domain.AdjustmentRequest, error) {
	var list []domain.AdjustmentRequest
	err := r.db.Preload("Product").Preload("Location").
		Where("requested_by = ?", userID).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *adjustmentRepository) FindAll() ([]domain.AdjustmentRequest, error) {
	var list []domain.AdjustmentRequest
	err := r.db.Preload("Product").Preload("Location").Preload("Requester").
		Preload("Approver").Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *adjustmentRepository) Update(a *domain.AdjustmentRequest) error {
	return r.db.Save(a).Error
}
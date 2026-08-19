package repository

import (
	"github.com/prasetyoary48/wms/internal/domain"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	Create(w *domain.Warehouse) error
	FindByID(id uint) (*domain.Warehouse, error)
	FindAll() ([]domain.Warehouse, error)
	Update(w *domain.Warehouse) error
	Delete(id uint) error
}

type warehouseRepository struct{ db *gorm.DB }

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) Create(w *domain.Warehouse) error { return r.db.Create(w).Error }

func (r *warehouseRepository) FindByID(id uint) (*domain.Warehouse, error) {
	var w domain.Warehouse
	if err := r.db.First(&w, id).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *warehouseRepository) FindAll() ([]domain.Warehouse, error) {
	var list []domain.Warehouse
	err := r.db.Order("id asc").Find(&list).Error
	return list, err
}

func (r *warehouseRepository) Update(w *domain.Warehouse) error { return r.db.Save(w).Error }
func (r *warehouseRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Warehouse{}, id).Error
}

// --- Location (rak) ---

type LocationRepository interface {
	Create(l *domain.Location) error
	FindByID(id uint) (*domain.Location, error)
	FindByWarehouse(warehouseID uint) ([]domain.Location, error)
	FindAll() ([]domain.Location, error)
	Update(l *domain.Location) error
	Delete(id uint) error
}

type locationRepository struct{ db *gorm.DB }

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(l *domain.Location) error { return r.db.Create(l).Error }

func (r *locationRepository) FindByID(id uint) (*domain.Location, error) {
	var l domain.Location
	if err := r.db.Preload("Warehouse").First(&l, id).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *locationRepository) FindByWarehouse(warehouseID uint) ([]domain.Location, error) {
	var list []domain.Location
	err := r.db.Where("warehouse_id = ?", warehouseID).Order("code asc").Find(&list).Error
	return list, err
}

func (r *locationRepository) FindAll() ([]domain.Location, error) {
	var list []domain.Location
	err := r.db.Preload("Warehouse").Order("id asc").Find(&list).Error
	return list, err
}

func (r *locationRepository) Update(l *domain.Location) error { return r.db.Save(l).Error }
func (r *locationRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Location{}, id).Error
}
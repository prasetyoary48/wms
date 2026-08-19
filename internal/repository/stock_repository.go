package repository

import (
	"github.com/prasetyoary48/wms/internal/domain"
	"gorm.io/gorm"
)

type StockRepository interface {
	// FindByProduct mengembalikan semua baris stok suatu produk di semua lokasi/rak —
	// inilah cara mengetahui di rak mana saja barang tersimpan.
	FindByProduct(productID uint) ([]domain.Stock, error)
	FindByLocation(locationID uint) ([]domain.Stock, error)
	FindOne(productID, locationID uint) (*domain.Stock, error)
	Upsert(tx *gorm.DB, productID, locationID uint, deltaQty int) error
	FindAll() ([]domain.Stock, error)
	WithTx(fn func(tx *gorm.DB) error) error
}

type stockRepository struct{ db *gorm.DB }

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) FindByProduct(productID uint) ([]domain.Stock, error) {
	var list []domain.Stock
	err := r.db.Preload("Location.Warehouse").Where("product_id = ?", productID).Find(&list).Error
	return list, err
}

func (r *stockRepository) FindByLocation(locationID uint) ([]domain.Stock, error) {
	var list []domain.Stock
	err := r.db.Preload("Product").Where("location_id = ?", locationID).Find(&list).Error
	return list, err
}

func (r *stockRepository) FindOne(productID, locationID uint) (*domain.Stock, error) {
	var s domain.Stock
	err := r.db.Where("product_id = ? AND location_id = ?", productID, locationID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Upsert menambah/mengurangi qty stok di suatu lokasi. Dipakai di dalam
// transaksi DB supaya perubahan stok selalu konsisten dengan StockMovement.
// deltaQty boleh negatif untuk pengurangan (stock out / transfer keluar).
func (r *stockRepository) Upsert(tx *gorm.DB, productID, locationID uint, deltaQty int) error {
	var s domain.Stock
	err := tx.Where("product_id = ? AND location_id = ?", productID, locationID).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		s = domain.Stock{ProductID: productID, LocationID: locationID, Quantity: deltaQty}
		return tx.Create(&s).Error
	}
	if err != nil {
		return err
	}
	return tx.Model(&s).Update("quantity", gorm.Expr("quantity + ?", deltaQty)).Error
}

func (r *stockRepository) FindAll() ([]domain.Stock, error) {
	var list []domain.Stock
	err := r.db.Preload("Product").Preload("Location.Warehouse").Find(&list).Error
	return list, err
}

func (r *stockRepository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
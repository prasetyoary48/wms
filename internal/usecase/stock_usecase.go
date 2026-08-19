package usecase

import (
	"errors"

	"github.com/prasetyoary48/wms/internal/domain"
	"github.com/prasetyoary48/wms/internal/repository"
	"gorm.io/gorm"
)

type StockUsecase interface {
	// StockIn menambah stok di suatu lokasi (barang masuk dari supplier), langsung berstatus done.
	StockIn(productID, locationID uint, qty int, note string, userID uint) (*domain.StockMovement, error)
	// StockOut mengurangi stok karena barang keluar/terjual, langsung berstatus done.
	StockOut(productID, locationID uint, qty int, note string, userID uint) (*domain.StockMovement, error)
	// LocationsOfProduct menjawab pertanyaan "barang ini ada di rak mana saja".
	LocationsOfProduct(productID uint) ([]domain.Stock, error)
	StockInLocation(locationID uint) ([]domain.Stock, error)
	AllStocks() ([]domain.Stock, error)
}

type stockUsecase struct {
	stockRepo repository.StockRepository
	moveRepo  repository.StockMovementRepository
}

func NewStockUsecase(stockRepo repository.StockRepository, moveRepo repository.StockMovementRepository) StockUsecase {
	return &stockUsecase{stockRepo: stockRepo, moveRepo: moveRepo}
}

func (u *stockUsecase) StockIn(productID, locationID uint, qty int, note string, userID uint) (*domain.StockMovement, error) {
	if qty <= 0 {
		return nil, errors.New("qty harus lebih besar dari 0")
	}

	var movement domain.StockMovement
	err := u.stockRepo.WithTx(func(tx *gorm.DB) error {
		if err := u.stockRepo.Upsert(tx, productID, locationID, qty); err != nil {
			return err
		}
		movement = domain.StockMovement{
			ProductID:    productID,
			ToLocationID: &locationID,
			Qty:          qty,
			Type:         domain.MovementTypeIn,
			Status:       domain.MovementStatusDone,
			Note:         note,
			CreatedBy:    userID,
		}
		return u.moveRepo.Create(tx, &movement)
	})
	if err != nil {
		return nil, err
	}
	return &movement, nil
}

func (u *stockUsecase) StockOut(productID, locationID uint, qty int, note string, userID uint) (*domain.StockMovement, error) {
	if qty <= 0 {
		return nil, errors.New("qty harus lebih besar dari 0")
	}

	current, err := u.stockRepo.FindOne(productID, locationID)
	if err != nil || current.Quantity < qty {
		return nil, errors.New("stok tidak mencukupi di lokasi ini")
	}

	var movement domain.StockMovement
	err = u.stockRepo.WithTx(func(tx *gorm.DB) error {
		if err := u.stockRepo.Upsert(tx, productID, locationID, -qty); err != nil {
			return err
		}
		movement = domain.StockMovement{
			ProductID:      productID,
			FromLocationID: &locationID,
			Qty:            qty,
			Type:           domain.MovementTypeOut,
			Status:         domain.MovementStatusDone,
			Note:           note,
			CreatedBy:      userID,
		}
		return u.moveRepo.Create(tx, &movement)
	})
	if err != nil {
		return nil, err
	}
	return &movement, nil
}

func (u *stockUsecase) LocationsOfProduct(productID uint) ([]domain.Stock, error) {
	return u.stockRepo.FindByProduct(productID)
}

func (u *stockUsecase) StockInLocation(locationID uint) ([]domain.Stock, error) {
	return u.stockRepo.FindByLocation(locationID)
}

func (u *stockUsecase) AllStocks() ([]domain.Stock, error) {
	return u.stockRepo.FindAll()
}
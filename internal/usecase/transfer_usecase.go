package usecase

import (
	"errors"

	"github.com/prasetyoary48/wms/internal/domain"
	"github.com/prasetyoary48/wms/internal/repository"
	"gorm.io/gorm"
)

type TransferUsecase interface {
	// RequestTransfer dipakai Staff untuk mengajukan pemindahan barang antar lokasi/gudang.
	// Movement dibuat dengan status "pending" dan belum mengubah Stock.
	RequestTransfer(productID, fromLocationID, toLocationID uint, qty int, note string, userID uint) (*domain.StockMovement, error)
	// ApproveTransfer dipakai SPV. Saat disetujui, baru Stock benar-benar berpindah.
	ApproveTransfer(movementID, approverID uint) error
	RejectTransfer(movementID, approverID uint, reason string) error
	PendingTransfers() ([]domain.StockMovement, error)
}

type transferUsecase struct {
	stockRepo repository.StockRepository
	moveRepo  repository.StockMovementRepository
}

func NewTransferUsecase(stockRepo repository.StockRepository, moveRepo repository.StockMovementRepository) TransferUsecase {
	return &transferUsecase{stockRepo: stockRepo, moveRepo: moveRepo}
}

func (u *transferUsecase) RequestTransfer(productID, fromLocationID, toLocationID uint, qty int, note string, userID uint) (*domain.StockMovement, error) {
	if qty <= 0 {
		return nil, errors.New("qty harus lebih besar dari 0")
	}
	if fromLocationID == toLocationID {
		return nil, errors.New("lokasi asal dan tujuan tidak boleh sama")
	}

	current, err := u.stockRepo.FindOne(productID, fromLocationID)
	if err != nil || current.Quantity < qty {
		return nil, errors.New("stok tidak mencukupi di lokasi asal")
	}

	movement := domain.StockMovement{
		ProductID:      productID,
		FromLocationID: &fromLocationID,
		ToLocationID:   &toLocationID,
		Qty:            qty,
		Type:           domain.MovementTypeTransfer,
		Status:         domain.MovementStatusPending,
		Note:           note,
		CreatedBy:      userID,
	}
	if err := u.moveRepo.Create(nil, &movement); err != nil {
		return nil, err
	}
	return &movement, nil
}

func (u *transferUsecase) ApproveTransfer(movementID, approverID uint) error {
	movement, err := u.moveRepo.FindByID(movementID)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if movement.Status != domain.MovementStatusPending {
		return errors.New("pengajuan sudah diproses sebelumnya")
	}

	return u.stockRepo.WithTx(func(tx *gorm.DB) error {
		if err := u.stockRepo.Upsert(tx, movement.ProductID, *movement.FromLocationID, -movement.Qty); err != nil {
			return err
		}
		if err := u.stockRepo.Upsert(tx, movement.ProductID, *movement.ToLocationID, movement.Qty); err != nil {
			return err
		}
		movement.Status = domain.MovementStatusApproved
		movement.ApprovedBy = &approverID
		return tx.Save(movement).Error
	})
}

func (u *transferUsecase) RejectTransfer(movementID, approverID uint, reason string) error {
	movement, err := u.moveRepo.FindByID(movementID)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if movement.Status != domain.MovementStatusPending {
		return errors.New("pengajuan sudah diproses sebelumnya")
	}
	movement.Status = domain.MovementStatusRejected
	movement.ApprovedBy = &approverID
	if reason != "" {
		movement.Note = movement.Note + " | ditolak: " + reason
	}
	return u.moveRepo.Update(movement)
}

func (u *transferUsecase) PendingTransfers() ([]domain.StockMovement, error) {
	return u.moveRepo.FindByStatus(domain.MovementStatusPending)
}
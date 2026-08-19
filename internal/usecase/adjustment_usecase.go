package usecase

import (
	"errors"

	"github.com/prasetyoary48/wms/internal/domain"
	"github.com/prasetyoary48/wms/internal/repository"
	"gorm.io/gorm"
)

type AdjustmentUsecase interface {
	Request(productID, locationID uint, requestedQty int, reason string, userID uint) (*domain.AdjustmentRequest, error)
	Approve(id, approverID uint) error
	Reject(id, approverID uint, reason string) error
	PendingList() ([]domain.AdjustmentRequest, error)
	MyRequests(userID uint) ([]domain.AdjustmentRequest, error)
}

type adjustmentUsecase struct {
	adjRepo   repository.AdjustmentRepository
	stockRepo repository.StockRepository
}

func NewAdjustmentUsecase(adjRepo repository.AdjustmentRepository, stockRepo repository.StockRepository) AdjustmentUsecase {
	return &adjustmentUsecase{adjRepo: adjRepo, stockRepo: stockRepo}
}

// Request dipakai Staff mengajukan koreksi jumlah stok (misal hasil temuan opname)
// di suatu lokasi. requestedQty adalah jumlah fisik yang seharusnya (bukan delta).
func (u *adjustmentUsecase) Request(productID, locationID uint, requestedQty int, reason string, userID uint) (*domain.AdjustmentRequest, error) {
	if requestedQty < 0 {
		return nil, errors.New("qty tidak boleh negatif")
	}
	adj := &domain.AdjustmentRequest{
		ProductID:    productID,
		LocationID:   locationID,
		Type:         domain.AdjustmentTypeStockFix,
		RequestedQty: requestedQty,
		Reason:       reason,
		Status:       domain.AdjustmentStatusPending,
		RequestedBy:  userID,
	}
	if err := u.adjRepo.Create(adj); err != nil {
		return nil, err
	}
	return adj, nil
}

func (u *adjustmentUsecase) Approve(id, approverID uint) error {
	adj, err := u.adjRepo.FindByID(id)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if adj.Status != domain.AdjustmentStatusPending {
		return errors.New("pengajuan sudah diproses sebelumnya")
	}

	current, err := u.stockRepo.FindOne(adj.ProductID, adj.LocationID)
	currentQty := 0
	if err == nil {
		currentQty = current.Quantity
	}
	delta := adj.RequestedQty - currentQty

	return u.stockRepo.WithTx(func(tx *gorm.DB) error {
		if delta != 0 {
			if err := u.stockRepo.Upsert(tx, adj.ProductID, adj.LocationID, delta); err != nil {
				return err
			}
		}
		adj.Status = domain.AdjustmentStatusApproved
		adj.ApprovedBy = &approverID
		return tx.Save(adj).Error
	})
}

func (u *adjustmentUsecase) Reject(id, approverID uint, reason string) error {
	adj, err := u.adjRepo.FindByID(id)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}
	if adj.Status != domain.AdjustmentStatusPending {
		return errors.New("pengajuan sudah diproses sebelumnya")
	}
	adj.Status = domain.AdjustmentStatusRejected
	adj.ApprovedBy = &approverID
	if reason != "" {
		adj.Reason = adj.Reason + " | ditolak: " + reason
	}
	return u.adjRepo.Update(adj)
}

func (u *adjustmentUsecase) PendingList() ([]domain.AdjustmentRequest, error) {
	return u.adjRepo.FindByStatus(domain.AdjustmentStatusPending)
}

func (u *adjustmentUsecase) MyRequests(userID uint) ([]domain.AdjustmentRequest, error) {
	return u.adjRepo.FindByUser(userID)
}
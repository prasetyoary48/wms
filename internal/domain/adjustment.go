package domain

import "time"

const (
	AdjustmentTypeStockFix = "stock_fix"
	AdjustmentTypeTransfer = "transfer"
	AdjustmentStatusPending  = "pending"
	AdjustmentStatusApproved = "approved"
	AdjustmentStatusRejected = "rejected"
)

// AdjustmentRequest adalah pengajuan dari Staff (misal koreksi stok atau
// permintaan pindah barang) yang wajib di-approve oleh SPV sebelum
// benar-benar mengubah data Stock.
type AdjustmentRequest struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProductID      uint      `gorm:"not null" json:"product_id"`
	Product        Product   `gorm:"foreignKey:ProductID" json:"product"`
	LocationID     uint      `gorm:"not null" json:"location_id"`
	Location       Location  `gorm:"foreignKey:LocationID" json:"location"`
	Type           string    `gorm:"size:20;not null" json:"type"`
	RequestedQty   int       `json:"requested_qty"`
	TargetLocation *uint     `json:"target_location_id"`
	Reason         string    `json:"reason"`
	Status         string    `gorm:"size:20;not null;default:pending" json:"status"`
	RequestedBy    uint      `gorm:"not null" json:"requested_by"`
	Requester      User      `gorm:"foreignKey:RequestedBy" json:"requester"`
	ApprovedBy     *uint     `json:"approved_by"`
	Approver       *User     `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (AdjustmentRequest) TableName() string { return "adjustment_requests" }
package domain

import "time"

// Stock adalah source of truth kuantitas produk pada suatu lokasi/rak.
// Setiap perubahan (in/out/transfer) harus mengubah baris ini melalui
// StockMovement supaya history tetap tercatat (audit trail).
type Stock struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `gorm:"not null;uniqueIndex:idx_product_location" json:"product_id"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"product"`
	LocationID uint      `gorm:"not null;uniqueIndex:idx_product_location" json:"location_id"`
	Location   Location  `gorm:"foreignKey:LocationID" json:"location"`
	Quantity   int       `gorm:"not null;default:0" json:"quantity"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Stock) TableName() string { return "stocks" }

const (
	MovementTypeIn       = "in"       // barang masuk dari supplier
	MovementTypeOut      = "out"      // barang keluar / terjual
	MovementTypeTransfer = "transfer" // pindah antar lokasi/gudang
)

const (
	MovementStatusPending  = "pending"
	MovementStatusApproved = "approved"
	MovementStatusRejected = "rejected"
	MovementStatusDone     = "done"
)

// StockMovement mencatat setiap pergerakan barang, termasuk transfer
// antar lokasi dalam satu gudang maupun antar gudang, atau keluar karena terjual.
type StockMovement struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProductID      uint      `gorm:"not null" json:"product_id"`
	Product        Product   `gorm:"foreignKey:ProductID" json:"product"`
	FromLocationID *uint     `json:"from_location_id"`
	FromLocation   *Location `gorm:"foreignKey:FromLocationID" json:"from_location,omitempty"`
	ToLocationID   *uint     `json:"to_location_id"`
	ToLocation     *Location `gorm:"foreignKey:ToLocationID" json:"to_location,omitempty"`
	Qty            int       `gorm:"not null" json:"qty"`
	Type           string    `gorm:"size:20;not null" json:"type"` // in / out / transfer
	Status         string    `gorm:"size:20;not null;default:pending" json:"status"`
	Note           string    `json:"note"`
	CreatedBy      uint      `gorm:"not null" json:"created_by"`
	Creator        User      `gorm:"foreignKey:CreatedBy" json:"creator"`
	ApprovedBy     *uint     `json:"approved_by"`
	Approver       *User     `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (StockMovement) TableName() string { return "stock_movements" }
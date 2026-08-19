package domain

import "time"

type Supplier struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:150;not null" json:"name"`
	Contact   string    `json:"contact"`
	Address   string    `json:"address"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Supplier) TableName() string { return "suppliers" }

const (
	TransactionTypeIn  = "in"
	TransactionTypeOut = "out"
)

// Transaction adalah dokumen header untuk penerimaan barang dari supplier
// atau pengeluaran barang karena terjual. Detail baris barang ada di
// StockMovement yang terhubung lewat Note/ReferenceNo (bisa dikembangkan
// jadi tabel transaction_details sesuai kebutuhan).
type Transaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ReferenceNo string    `gorm:"uniqueIndex;size:50;not null" json:"reference_no"`
	Type        string    `gorm:"size:20;not null" json:"type"`
	SupplierID  *uint     `json:"supplier_id"`
	Supplier    *Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	WarehouseID uint      `gorm:"not null" json:"warehouse_id"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	Creator     User      `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Transaction) TableName() string { return "transactions" }
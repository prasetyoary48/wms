package domain

import "time"

const (
	OpnameStatusOngoing  = "ongoing"
	OpnameStatusFinished = "finished"
	OpnameStatusApproved = "approved"
)

// StockOpname adalah sesi audit stok fisik di satu gudang.
type StockOpname struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WarehouseID uint      `gorm:"not null" json:"warehouse_id"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse"`
	Status      string    `gorm:"size:20;not null;default:ongoing" json:"status"`
	ConductedBy uint      `gorm:"not null" json:"conducted_by"`
	Conductor   User      `gorm:"foreignKey:ConductedBy" json:"conductor"`
	ApprovedBy  *uint     `json:"approved_by"`
	Approver    *User     `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (StockOpname) TableName() string { return "stock_opname" }

// StockOpnameDetail mencatat perbandingan qty sistem vs qty fisik per produk per lokasi.
// Jika ada selisih, biasanya akan otomatis dibuatkan AdjustmentRequest.
type StockOpnameDetail struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	OpnameID   uint     `gorm:"not null" json:"opname_id"`
	ProductID  uint     `gorm:"not null" json:"product_id"`
	Product    Product  `gorm:"foreignKey:ProductID" json:"product"`
	LocationID uint     `gorm:"not null" json:"location_id"`
	Location   Location `gorm:"foreignKey:LocationID" json:"location"`
	SystemQty  int      `gorm:"not null" json:"system_qty"`
	ActualQty  int      `gorm:"not null" json:"actual_qty"`
	Note       string   `json:"note"`
}

func (StockOpnameDetail) TableName() string { return "stock_opname_details" }
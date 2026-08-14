package domain

import "time"

type Warehouse struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"uniqueIndex;size:20;not null" json:"code"`
	Name      string    `gorm:"size:150;not null" json:"name"`
	Address   string    `json:"address"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Warehouse) TableName() string { return "warehouses" }

// Location merepresentasikan rak/lokasi penyimpanan di dalam gudang.
// Code disarankan berformat hierarkis, misal "A-01-02" (Zona A, Rak 01, Level 02),
// supaya mudah digenerate jadi barcode/QR untuk scanning.
type Location struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WarehouseID uint      `gorm:"not null" json:"warehouse_id"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse"`
	Code        string    `gorm:"size:30;not null" json:"code"`
	Zone        string    `gorm:"size:20" json:"zone"`
	Capacity    int       `json:"capacity"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Location) TableName() string { return "locations" }
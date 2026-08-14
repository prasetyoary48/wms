package domain

import "time"

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;not null" json:"name"`
}

func (Category) TableName() string { return "categories" }

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SKU         string    `gorm:"uniqueIndex;size:50;not null" json:"sku"`
	Name        string    `gorm:"size:150;not null" json:"name"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Unit        string    `gorm:"size:20;not null" json:"unit"` // pcs, box, kg, dll
	Description string    `json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Product) TableName() string { return "products" }
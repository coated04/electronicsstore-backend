package models

import brandmodels "device-store/models" 

type Device struct {
	ID         uint                 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string               `json:"name"`
	Type       string               `json:"type"`
	BrandID    uint                 `json:"brand_id"`
	Brand      brandmodels.Brand    `gorm:"foreignKey:BrandID"`
	Price      float64          `json:"price"`
	CategoryID uint                 `json:"category_id"`
	Category   brandmodels.Category `gorm:"foreignKey:CategoryID"`
	
}



type Brand struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null;unique" json:"name"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null;unique" json:"name"`
}

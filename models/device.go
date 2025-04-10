package models

type Device struct {
    ID         uint    `gorm:"primaryKey" json:"id"`
    Name       string  `json:"name"`
    BrandID    uint    `json:"brand_id"`
    Brand      Brand   `gorm:"foreignKey:BrandID"`
    CategoryID uint    `json:"category_id"`
    Category   Category `gorm:"foreignKey:CategoryID"`
    Price      float64 `json:"price"`
}

package models

type Brand struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `json:"name"`
}

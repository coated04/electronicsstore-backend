package models

type User struct {
    ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Username string `gorm:"unique;not null" json:"username"`
    Email    string `gorm:"unique;not null" json:"email"`
    Password string `gorm:"not null" json:"password"`
}
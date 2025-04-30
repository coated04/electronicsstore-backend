package models

import usermodels "device-store/user-service/models"

type Device struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string             `json:"name"`
	UserID uint               `json:"user_id"`
	Type   string `json:"type"`
	User   usermodels.User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

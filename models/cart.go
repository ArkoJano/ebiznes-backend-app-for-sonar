package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	Items []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	CartID    uint    `json:"cart_id"`
	Quantity  int     `json:"quantity"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	Amount  float64 `json:"amount"`
	OrderID uint    `json:"order_id"`
}

package utils

import "time"

type Order struct {
	ID        uint      `gorm:"column:id"`
	OrderUUID string    `gorm:"column:orderUUID"`
	Symbol    string    `gorm:"column:symbol"`
	Volume    float32   `gorm:"column:volume"`
	Price     float64   `gorm:"column:price"`
	Type      string    `gorm:"column:type"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type OrderResponse struct {
	OrderID   uint      `json:"orderId"`
	OrderUUID string    `json:"orderUUID"`
	Symbol    string    `json:"symbol"`
	Volume    float32   `json:"volume"`
	Price     float64   `json:"price"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

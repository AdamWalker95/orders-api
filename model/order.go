package model

import (
	"time"
)

type Order struct {
	OrderID     uint64     `json:"order_id"`
	CustomerID  uint64     `json:"customer_id"`
	CreatedAt   *time.Time `json:"created_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

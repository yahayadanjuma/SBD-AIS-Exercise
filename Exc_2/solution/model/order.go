package model

import "time"

// Order represents a customer's drink order.
type Order struct {
	ID        uint64    `json:"id"`
	DrinkID   uint64    `json:"drink_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

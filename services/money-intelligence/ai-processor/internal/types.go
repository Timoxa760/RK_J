package internal

import "time"

type RawReceipt struct {
	ID       string    `json:"id,omitempty"`
	UserID   string    `json:"user_id,omitempty"`
	Provider string    `json:"provider"`
	Store    string    `json:"store_name"`
	Date     time.Time `json:"purchased_at"`
	Total    float64   `json:"total_amount"`
	Items    []Item    `json:"items"`
}

type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type CategorizedItem struct {
	Name       string
	Price      float64
	Quantity   int
	Category   string
	IsImpulsive bool
}

type CategorizedReceipt struct {
	UserID   string
	Store    string
	Provider string
	Date     time.Time
	Items    []CategorizedItem
}

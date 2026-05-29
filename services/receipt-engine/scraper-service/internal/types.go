package internal

import "time"

type RawItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type RawReceipt struct {
	ID       string    `json:"id,omitempty"`
	UserID   string    `json:"user_id,omitempty"`
	Provider string    `json:"provider"`
	Store    string    `json:"store_name"`
	Date     time.Time `json:"purchased_at"`
	Total    float64   `json:"total_amount"`
	Items    []RawItem `json:"items"`
}

package order

import "time"

type Order struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Items  []Item
}

type Item struct {
	ID       int     `json:"id"`
	Quantity int     `json:"quantity"`
	Type     string  `json:"type"`
	Product  Product `json:"product"`
}
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	SkuId       string    `json:"skuId"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

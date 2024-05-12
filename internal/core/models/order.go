package models

type ProductionOrderPage struct {
	Results []ProductionOrder `json:"results"`
	Next    *int              `json:"next"`
}

type ProductionOrder struct {
	ID       int       `json:"id"`
	Status   string    `json:"status"`
	Products []Product `json:"products"`
}

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

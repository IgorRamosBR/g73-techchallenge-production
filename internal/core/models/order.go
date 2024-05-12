package models

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

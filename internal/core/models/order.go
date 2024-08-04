package models

import "time"

type Order struct {
	ID          string      `json:"id" dynamodbav:"PK"`
	Status      string      `json:"status" dynamodbav:"Status"`
	CustomerCPF string      `json:"customerCPF" dynamodbav:"CustomerCPF"`
	CreatedAt   time.Time   `json:"createdAt" dynamodbav:"CreatedAt"`
	FinishedAt  time.Time   `json:"updatedAt" dynamodbav:"FinishedAt"`
	Items       []OrderItem `json:"items" dynamodbav:"Items"`
	Entity      string      `json:"entity" dynamodbav:"GSI1PK"`
}

type OrderItem struct {
	Quantity int     `json:"quantity" dynamodbav:"Quantity"`
	Type     string  `json:"type" dynamodbav:"Type"`
	Product  Product `json:"product" dynamodbav:"Product"`
}

type Product struct {
	Name        string `json:"name" dynamodbav:"Name"`
	Description string `json:"description" dynamodbav:"Description"`
}

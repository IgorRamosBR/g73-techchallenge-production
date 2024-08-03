package gateways

import (
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-payment/pkg/dynamodb"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type OrderRepository interface {
	GetOrders() ([]models.Order, error)
	SaveOrder(order models.Order) error
	UpdateOrderStatus(orderId int, status string) error
}

type orderRepository struct {
	table          string
	dynamodbClient dynamodb.DynamoDBClient
}

func NewOrderRepository(dynamodbClient dynamodb.DynamoDBClient) OrderRepository {
	return &orderRepository{
		dynamodbClient: dynamodbClient,
	}
}

func (r *orderRepository) GetOrders() ([]models.Order, error) {
	type Key struct {
		entity string `dynamodbav:"entity"`
	}
	key := Key{entity: "ORDER"}
	keyMap, err := attributevalue.MarshalMap(key)
	if err != nil {
		return nil, fmt.Errorf("failed to map key: %w", err)
	}

	items, err := r.dynamodbClient.GetItem(r.table, keyMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := []models.Order{}
	err = attributevalue.UnmarshalMap(items, &orders)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal orders: %w", err)
	}

	return orders, nil
}

func (r *orderRepository) SaveOrder(order models.Order) error {
	av, err := attributevalue.MarshalMap(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	err = r.dynamodbClient.PutItem(r.table, av)
	if err != nil {
		return fmt.Errorf("failed to put order: %w", err)
	}

	return nil
}

func (r *orderRepository) UpdateOrderStatus(orderId int, status string) error {
	id, err := attributevalue.Marshal(orderId)
	if err != nil {
		return fmt.Errorf("failed to marshal order id: %w", err)
	}

	key := map[string]types.AttributeValue{"OrderID": id}

	update := expression.Set(expression.Name("Status"), expression.Value(status))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to create expression: %w", err)
	}

	err = r.dynamodbClient.UpdateItem(r.table, key, expr)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

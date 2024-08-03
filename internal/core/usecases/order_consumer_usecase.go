package usecases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/drivers/broker"
)

type OrderConsumerUseCase interface {
	StartConsumers()
}

type orderConsumerUseCase struct {
	orderPaidConsumer broker.Consumer
	orderPublisher    broker.Publisher
	orderUsecase      OrderUseCase
}

type OrderConsumerUseCaseConfig struct {
	OrderPaidConsumer broker.Consumer
	OrderPublisher    broker.Publisher
	OrderUseCase      OrderUseCase
}

func NewOrderConsumerUseCase(orderPaidConsumer broker.Consumer, orderPublisher broker.Publisher, orderUsecase OrderUseCase) OrderConsumerUseCase {
	return &orderConsumerUseCase{
		orderPaidConsumer: orderPaidConsumer,
		orderPublisher:    orderPublisher,
		orderUsecase:      orderUsecase,
	}
}

func (u *orderConsumerUseCase) StartConsumers() {
	go u.orderPaidConsumer.StartConsumer(u.processOrderMessage)
}

func (u *orderConsumerUseCase) processOrderMessage(message []byte) error {
	var productionOrder events.OrderProductionDTO
	err := json.Unmarshal(message, &productionOrder)
	if err != nil {
		return fmt.Errorf("failed to unmarshall message, error: %w", err)
	}

	order := mapEventOrderToOrder(productionOrder)
	err = u.orderUsecase.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("failed to update order status, error: %w", err)
	}

	return nil
}

func mapEventOrderToOrder(productionOrder events.OrderProductionDTO) models.Order {
	orderItems := make([]models.OrderItem, len(productionOrder.Items))
	for i, item := range productionOrder.Items {
		orderItem := models.OrderItem{
			Quantity: item.Quantity,
			Type:     item.Type,
			Product: models.Product{
				Name:        item.Products.Name,
				Description: item.Products.Description,
			},
		}
		orderItems[i] = orderItem
	}

	order := models.Order{
		ID:        productionOrder.ID,
		Status:    "CREATED",
		CreatedAt: time.Now(),
		Items:     orderItems,
		Entity:    "ORDER",
	}

	return order
}

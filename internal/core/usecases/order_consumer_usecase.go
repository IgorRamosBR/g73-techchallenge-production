package usecases

import (
	"encoding/json"
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events"
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

	err = u.orderUsecase.UpdateOrderStatus(productionOrder.ID, productionOrder.Status)
	if err != nil {
		return fmt.Errorf("failed to update order status, error: %w", err)
	}

	return nil
}

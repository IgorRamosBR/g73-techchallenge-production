package gateways

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/drivers/broker"
)

type OrderNotify interface {
	NotifyOrder(orderId int, status string) error
}

type orderNotify struct {
	publisher   broker.Publisher
	destination string
}

func NewOrderNotify(publisher broker.Publisher, destination string) OrderNotify {
	return orderNotify{publisher: publisher, destination: destination}
}

func (o orderNotify) NotifyOrder(orderId int, status string) error {
	message, err := json.Marshal(events.OrderStatusEventDTO{
		OrderId: orderId,
		Status:  status,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payment order[%d] with status[%s], error: %v", orderId, status, err)
	}

	ctx := context.Background()
	err = o.publisher.Publish(ctx, o.destination, message)
	if err != nil {
		return fmt.Errorf("failed to publish order[%d] with status[%s], error: %v", orderId, status, err)
	}

	return nil
}

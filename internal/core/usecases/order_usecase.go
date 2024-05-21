package usecases

import (
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	gateways "github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/gateways/order"
)

type OrderUseCase interface {
	GetOrders() (models.ProductionOrderPage, error)
	UpdateOrderStatus(orderId string, status string) error
}

type orderUseCase struct {
	orderClient gateways.OrderClient
}

func NewOrderUseCase(orderClient gateways.OrderClient) OrderUseCase {
	return orderUseCase{
		orderClient: orderClient,
	}
}

func (o orderUseCase) GetOrders() (models.ProductionOrderPage, error) {
	orders, err := o.orderClient.GetOrders()
	if err != nil {
		return models.ProductionOrderPage{}, err
	}

	return orders, nil
}

func (o orderUseCase) UpdateOrderStatus(orderId string, status string) error {
	err := o.orderClient.UpdateOrderStatus(orderId, status)
	if err != nil {
		return err
	}

	return nil
}

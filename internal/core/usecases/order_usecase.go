package usecases

import (
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/gateways"
)

type OrderUseCase interface {
	GetOrders() ([]models.Order, error)
	CreateOrder(order models.Order) error
	UpdateOrderStatus(orderId int, orderStatus string) error
}

type orderUseCase struct {
	orderRepository gateways.OrderRepository
	orderNotify     gateways.OrderNotify
}

func NewOrderUseCase(orderRepository gateways.OrderRepository, orderNotify gateways.OrderNotify) OrderUseCase {
	return &orderUseCase{
		orderRepository: orderRepository,
		orderNotify:     orderNotify,
	}
}

func (o orderUseCase) GetOrders() ([]models.Order, error) {
	orders, err := o.orderRepository.GetOrders()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderUseCase) UpdateOrderStatus(orderId int, orderStatus string) error {
	err := o.orderRepository.UpdateOrderStatus(orderId, orderStatus)
	if err != nil {
		return err
	}

	return nil
}

func (o *orderUseCase) CreateOrder(order models.Order) error {
	err := o.orderRepository.SaveOrder(order)
	if err != nil {
		return err
	}

	return nil
}

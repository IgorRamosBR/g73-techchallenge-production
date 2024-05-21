package usecases

import (
	"errors"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	mock_order "github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/gateways/order/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderUseCase_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderClient := mock_order.NewMockOrderClient(ctrl)

	orderUseCase := NewOrderUseCase(orderClient)

	orderClient.EXPECT().GetOrders().Times(1).Return(models.ProductionOrderPage{}, errors.New("internal error"))

	orders, err := orderUseCase.GetOrders()

	assert.Empty(t, orders)
	assert.EqualError(t, err, "internal error")

	expectedOrders := models.ProductionOrderPage{
		Results: []models.ProductionOrder{createOrder()},
		Next:    new(int),
	}
	orderClient.EXPECT().GetOrders().Times(1).Return(expectedOrders, nil)

	orders, err = orderUseCase.GetOrders()

	assert.Equal(t, expectedOrders, orders)
	assert.NoError(t, err)
}

func TestOrderUseCase_UpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	orderClient := mock_order.NewMockOrderClient(ctrl)

	orderUseCase := NewOrderUseCase(orderClient)

	orderClient.EXPECT().UpdateOrderStatus("123", "PAID").Times(1).Return(errors.New("internal error"))

	err := orderUseCase.UpdateOrderStatus("123", "PAID")

	assert.EqualError(t, err, "internal error")

	orderClient.EXPECT().UpdateOrderStatus("123", "PAID").Times(1).Return(nil)

	err = orderUseCase.UpdateOrderStatus("123", "PAID")

	assert.NoError(t, err)
}

func createOrder() models.ProductionOrder {
	return models.ProductionOrder{
		ID:     123,
		Status: "PAID",
		Products: []models.Product{
			{
				Name:        "Batata Frita",
				Description: "Batata canoa",
				Category:    "Acompanhamento",
			},
		},
	}
}

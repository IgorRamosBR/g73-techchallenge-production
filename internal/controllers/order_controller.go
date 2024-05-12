package controllers

import (
	"net/http"

	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/usecases"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/usecases/dto"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderUseCase usecases.OrderUseCase
}

func NewOrderController(orderUseCase usecases.OrderUseCase) OrderController {
	return OrderController{
		orderUseCase: orderUseCase,
	}
}

func (o OrderController) GetOrdersHandler(c *gin.Context) {
	orders, err := o.orderUseCase.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (o OrderController) UpdateOrderStatusHandler(c *gin.Context) {
	orderId := c.Param("id")

	var orderStatusRequest dto.OrderStatusRequest
	err := c.BindJSON(&orderStatusRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requests"})
		return
	}

	err = o.orderUseCase.UpdateOrderStatus(orderId, orderStatusRequest.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update order status"})
		return
	}

	c.Status(http.StatusNoContent)
}

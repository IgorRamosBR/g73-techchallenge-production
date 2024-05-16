package api

import (
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/controllers"
	"github.com/gin-gonic/gin"
)

func NewApi(orderController controllers.OrderController) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/orders", orderController.GetOrdersHandler)
		v1.PUT("/orders/:id/status", orderController.UpdateOrderStatusHandler)
	}

	return router
}

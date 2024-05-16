package main

import (
	"github.com/IgorRamosBR/g73-techchallenge-production/configs"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/api"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/controllers"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/usecases"
	gateways "github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/gateways/order"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/http"
)

func main() {
	appConfig := configs.GetAppConfig()

	httpClient := http.NewHttpClient(appConfig.OrderTimeout)
	orderClient := gateways.NewOrderClient(httpClient, appConfig.OrderUrl)

	orderUseCase := usecases.NewOrderUseCase(orderClient)
	orderController := controllers.NewOrderController(orderUseCase)

	api := api.NewApi(orderController)
	api.Run(":" + appConfig.Port)
}

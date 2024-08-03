package configs

import (
	"os"
)

type AppConfig struct {
	Port string

	OrderTable         string
	OrderTableEndpoint string

	OrderEventsBrokerUrl        string
	OrderPaidEventsTopic        string
	OrderReadyEventsDestination string
}

func GetAppConfig() AppConfig {
	appConfig := AppConfig{}

	appConfig.Port = os.Getenv("PORT")
	appConfig.OrderTable = os.Getenv("ORDER_TABLE")
	appConfig.OrderTableEndpoint = os.Getenv("ORDER_TABLE_ENDPOINT")
	appConfig.OrderEventsBrokerUrl = os.Getenv("ORDER_EVENTS_BROKER_URL")
	appConfig.OrderPaidEventsTopic = os.Getenv("ORDER_PAID_EVENTS")
	appConfig.OrderReadyEventsDestination = os.Getenv("ORDER_READY_EVENTS_DESTINATION")

	return appConfig
}

package configs

import (
	"os"
	"time"
)

type AppConfig struct {
	OrderUrl     string
	OrderTimeout time.Duration
}

func GetAppConfig() AppConfig {
	appConfig := AppConfig{}

	appConfig.OrderUrl = os.Getenv("ORDER_URL")
	orderTimeout := os.Getenv("ORDER_TIMEOUT_MS")
	orderTimeoutTime, err := time.ParseDuration(orderTimeout)
	if err != nil {
		panic(err)
	}

	appConfig.OrderTimeout = orderTimeoutTime
	return appConfig
}

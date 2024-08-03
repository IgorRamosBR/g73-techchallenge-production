package main

import (
	"context"

	"github.com/IgorRamosBR/g73-techchallenge-payment/pkg/dynamodb"
	"github.com/IgorRamosBR/g73-techchallenge-production/configs"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/api"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/controllers"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/usecases"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/drivers/broker"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/gateways"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsDynamoDb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	appConfig := configs.GetAppConfig()

	dynamodbClient, err := NewDynamoDBClient(appConfig.OrderTableEndpoint)
	if err != nil {
		panic(err)
	}

	publisher, err := NewRabbitMQPublisher(appConfig.OrderEventsBrokerUrl, appConfig.OrderReadyEventsDestination)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	orderRepository := gateways.NewOrderRepository(dynamodbClient)
	orderNotify := gateways.NewOrderNotify(publisher, appConfig.OrderPaidEventsTopic)
	orderUseCase := usecases.NewOrderUseCase(orderRepository, orderNotify)
	orderController := controllers.NewOrderController(orderUseCase)

	api := api.NewApi(orderController)
	api.Run(":" + appConfig.Port)
}

func NewDynamoDBClient(endpoint string) (dynamodb.DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	if endpoint != "" {
		client := awsDynamoDb.NewFromConfig(cfg, func(o *awsDynamoDb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
		return dynamodb.NewDynamoDBClient(client), nil
	}

	client := awsDynamoDb.NewFromConfig(cfg)
	return dynamodb.NewDynamoDBClient(client), nil

}

func NewRabbitMQPublisher(url, exchange string) (broker.Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	publisher := broker.NewRabbitMQPublisher(conn, ch, exchange)

	return publisher, nil
}

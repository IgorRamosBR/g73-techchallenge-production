package main

import (
	"context"

	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events/broker"
	"github.com/IgorRamosBR/g73-techchallenge-payment/pkg/dynamodb"
	"github.com/IgorRamosBR/g73-techchallenge-production/configs"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/api"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/controllers"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/usecases"
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

	brokerChannel, err := NewRabbitMQBrokerChannel(appConfig.OrderEventsBrokerUrl)
	if err != nil {
		panic(err)
	}
	defer brokerChannel.Close()

	ordersPaidQueue, err := broker.NewRabbitMQConsumer(brokerChannel, appConfig.OrderInProgressEventsQueue)
	if err != nil {
		panic(err)
	}

	publisher := broker.NewRabbitMQPublisher(brokerChannel, appConfig.OrderEventsTopic)
	defer publisher.Close()

	orderRepository := gateways.NewOrderRepository(dynamodbClient, appConfig.OrderTable)
	orderNotify := gateways.NewOrderNotify(publisher, appConfig.OrderReadyEventsDestination)
	orderUseCase := usecases.NewOrderUseCase(orderRepository, orderNotify)
	orderConsumerUseCase := usecases.NewOrderConsumerUseCase(ordersPaidQueue, orderUseCase)
	orderConsumerUseCase.StartConsumers()

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

func NewRabbitMQBrokerChannel(url string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, err
}

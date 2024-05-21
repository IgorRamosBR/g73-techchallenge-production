package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/drivers/http"
)

type OrderClient interface {
	GetOrders() (models.ProductionOrderPage, error)
	UpdateOrderStatus(orderId string, status string) error
}

type orderClient struct {
	httpClient http.HttpClient
	orderUrl   string
}

func NewOrderClient(httpClient http.HttpClient, orderUrl string) OrderClient {
	return &orderClient{
		httpClient: httpClient,
		orderUrl:   orderUrl,
	}
}

func (c *orderClient) GetOrders() (models.ProductionOrderPage, error) {
	resp, err := c.httpClient.DoGet(c.orderUrl)
	if err != nil {
		return models.ProductionOrderPage{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return models.ProductionOrderPage{}, errors.New("resp status code non-2xx")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ProductionOrderPage{}, err
	}

	var ordersPage OrderPage
	if err := json.Unmarshal(body, &ordersPage); err != nil {
		return models.ProductionOrderPage{}, err
	}

	productionOrders := c.mapOrdersToProductionOrders(ordersPage.Results)

	return models.ProductionOrderPage{
		Results: productionOrders,
		Next:    new(int),
	}, nil
}

func (c *orderClient) UpdateOrderStatus(orderId string, status string) error {
	orderStatusRequest := dto.OrderStatusRequest{
		Status: status,
	}

	reqBody, err := json.Marshal(orderStatusRequest)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.DoPut(fmt.Sprintf("%s/%s/status", c.orderUrl, orderId), reqBody)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("resp status code non-2xx")
	}

	return nil
}

func (c *orderClient) mapOrdersToProductionOrders(orders []Order) []models.ProductionOrder {
	productionOrders := []models.ProductionOrder{}
	for _, o := range orders {
		products := c.mapProducts(o.Items)
		order := models.ProductionOrder{
			ID:       o.ID,
			Status:   o.Status,
			Products: products,
		}
		productionOrders = append(productionOrders, order)
	}

	return productionOrders
}

func (c *orderClient) mapProducts(items []Item) []models.Product {
	productionProducts := []models.Product{}
	for _, item := range items {
		product := models.Product{
			Name:        item.Product.Name,
			Description: item.Product.Description,
			Category:    item.Product.Category,
		}
		productionProducts = append(productionProducts, product)
	}
	return productionProducts
}

package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/http"
)

type OrderClient interface {
	GetOrders() ([]models.ProductionOrder, error)
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

func (c *orderClient) GetOrders() ([]models.ProductionOrder, error) {
	resp, err := c.httpClient.DoGet(c.orderUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("resp status code non-2xx")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(body, &orders); err != nil {
		return nil, err
	}

	return c.mapOrdersToProductionOrders(orders), nil
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

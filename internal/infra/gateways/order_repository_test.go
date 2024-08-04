package gateways

import (
	"errors"
	"testing"

	mock_dynamodb "github.com/IgorRamosBR/g73-techchallenge-payment/pkg/dynamodb/mocks"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	table := "Kitchen"
	gsi := "SecondaryIndex"
	mockDynamoDBClient := mock_dynamodb.NewMockDynamoDBClient(ctrl)
	repo := NewOrderRepository(mockDynamoDBClient, table)

	tests := []struct {
		name           string
		mockSetup      func()
		expectedOrders []models.Order
		expectedError  error
	}{
		{
			name: "success",
			mockSetup: func() {
				expectedOrders := []models.Order{
					{ID: "1", Status: "READY"},
					{ID: "2", Status: "READY"},
				}
				var expectedItems []map[string]types.AttributeValue
				for _, order := range expectedOrders {
					item, _ := attributevalue.MarshalMap(order)
					expectedItems = append(expectedItems, item)
				}

				entityExpr := expression.Key("GSI1PK").Equal(expression.Value("ORDER"))
				expr, _ := expression.NewBuilder().WithKeyCondition(entityExpr).Build()

				mockDynamoDBClient.EXPECT().
					QueryItem(table, expr, gsi).
					Return(expectedItems, nil)
			},
			expectedOrders: []models.Order{
				{ID: "1", Status: "READY"},
				{ID: "2", Status: "READY"},
			},
			expectedError: nil,
		},
		{
			name: "dynamodb error",
			mockSetup: func() {
				entityExpr := expression.Key("GSI1PK").Equal(expression.Value("ORDER"))
				expr, _ := expression.NewBuilder().WithKeyCondition(entityExpr).Build()

				mockDynamoDBClient.EXPECT().
					QueryItem(table, expr, gsi).
					Return(nil, errors.New("dynamodb error"))
			},
			expectedOrders: nil,
			expectedError:  errors.New("failed to get orders: dynamodb error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			orders, err := repo.GetOrders()
			assert.Equal(t, len(tt.expectedOrders), len(orders))
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	table := "Kitchen"
	mockDynamoDBClient := mock_dynamodb.NewMockDynamoDBClient(ctrl)
	repo := NewOrderRepository(mockDynamoDBClient, table)

	tests := []struct {
		name          string
		order         models.Order
		mockSetup     func()
		expectedError error
	}{
		{
			name:  "success",
			order: models.Order{ID: "1", Status: "NEW"},
			mockSetup: func() {
				order := models.Order{ID: "1", Status: "NEW"}
				orderAV, _ := attributevalue.MarshalMap(order)
				mockDynamoDBClient.EXPECT().
					PutItem(table, orderAV).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "dynamodb error",
			order: models.Order{ID: "1", Status: "NEW"},
			mockSetup: func() {
				order := models.Order{ID: "1", Status: "NEW"}
				orderAV, _ := attributevalue.MarshalMap(order)
				mockDynamoDBClient.EXPECT().
					PutItem(table, orderAV).
					Return(errors.New("dynamodb error"))
			},
			expectedError: errors.New("failed to put order: dynamodb error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.SaveOrder(tt.order)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	table := "Kitchen"
	mockDynamoDBClient := mock_dynamodb.NewMockDynamoDBClient(ctrl)
	repo := NewOrderRepository(mockDynamoDBClient, table)

	tests := []struct {
		name          string
		ID            int
		status        string
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "success",
			ID:     1,
			status: "COMPLETED",
			mockSetup: func() {
				ID := 1
				status := "COMPLETED"
				IDAV, _ := attributevalue.Marshal(ID)
				key := map[string]types.AttributeValue{"OrderID": IDAV}

				update := expression.Set(expression.Name("Status"), expression.Value(status))
				expr, _ := expression.NewBuilder().WithUpdate(update).Build()

				mockDynamoDBClient.EXPECT().
					UpdateItem(table, key, expr).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "dynamodb error",
			ID:     1,
			status: "COMPLETED",
			mockSetup: func() {
				ID := 1
				status := "COMPLETED"
				IDAV, _ := attributevalue.Marshal(ID)
				key := map[string]types.AttributeValue{"OrderID": IDAV}

				update := expression.Set(expression.Name("Status"), expression.Value(status))
				expr, _ := expression.NewBuilder().WithUpdate(update).Build()

				mockDynamoDBClient.EXPECT().
					UpdateItem(table, key, expr).
					Return(errors.New("dynamodb error"))
			},
			expectedError: errors.New("failed to update order status: dynamodb error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.UpdateOrderStatus(tt.ID, tt.status)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

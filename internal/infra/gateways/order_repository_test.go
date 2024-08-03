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

	mockDynamoDBClient := mock_dynamodb.NewMockDynamoDBClient(ctrl)
	repo := NewOrderRepository(mockDynamoDBClient)

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
					{ID: 1, Status: "NEW"},
					{ID: 2, Status: "PROCESSING"},
				}
				expectedItems, _ := attributevalue.MarshalMap(expectedOrders)
				mockDynamoDBClient.EXPECT().
					GetItem(gomock.Any(), gomock.Any()).
					Return(expectedItems, nil)
			},
			expectedOrders: []models.Order{
				{ID: 1, Status: "NEW"},
				{ID: 2, Status: "PROCESSING"},
			},
			expectedError: nil,
		},
		{
			name: "dynamodb error",
			mockSetup: func() {
				mockDynamoDBClient.EXPECT().
					GetItem(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("dynamodb error"))
			},
			expectedOrders: nil,
			expectedError:  errors.New("failed to get orders: dynamodb error"),
		},
		{
			name: "unmarshal error",
			mockSetup: func() {
				mockDynamoDBClient.EXPECT().
					GetItem(gomock.Any(), gomock.Any()).
					Return(map[string]types.AttributeValue{}, nil)
			},
			expectedOrders: nil,
			expectedError:  errors.New("failed to unmarshal orders: unmarshal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			orders, err := repo.GetOrders()
			assert.Equal(t, tt.expectedOrders, orders)
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

	mockDynamoDBClient := mock_dynamodb.NewMockDynamoDBClient(ctrl)
	repo := NewOrderRepository(mockDynamoDBClient)

	tests := []struct {
		name          string
		order         models.Order
		mockSetup     func()
		expectedError error
	}{
		{
			name:  "success",
			order: models.Order{ID: 1, Status: "NEW"},
			mockSetup: func() {
				order := models.Order{ID: 1, Status: "NEW"}
				orderAV, _ := attributevalue.MarshalMap(order)
				mockDynamoDBClient.EXPECT().
					PutItem(gomock.Any(), orderAV).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "dynamodb error",
			order: models.Order{ID: 1, Status: "NEW"},
			mockSetup: func() {
				order := models.Order{ID: 1, Status: "NEW"}
				orderAV, _ := attributevalue.MarshalMap(order)
				mockDynamoDBClient.EXPECT().
					PutItem(gomock.Any(), orderAV).
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

	mockDynamoDBClient := mock_dynamodb.NewMockDynamoDBClient(ctrl)
	repo := NewOrderRepository(mockDynamoDBClient)

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
				key := map[string]types.AttributeValue{"ID": IDAV}

				update := expression.Set(expression.Name("Status"), expression.Value(status))
				expr, _ := expression.NewBuilder().WithUpdate(update).Build()

				mockDynamoDBClient.EXPECT().
					UpdateItem(gomock.Any(), key, expr).
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
				key := map[string]types.AttributeValue{"ID": IDAV}

				update := expression.Set(expression.Name("Status"), expression.Value(status))
				expr, _ := expression.NewBuilder().WithUpdate(update).Build()

				mockDynamoDBClient.EXPECT().
					UpdateItem(gomock.Any(), key, expr).
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

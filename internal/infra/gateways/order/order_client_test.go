package order

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	mock_http "github.com/IgorRamosBR/g73-techchallenge-production/internal/infra/drivers/http/mocks"
	"go.uber.org/mock/gomock"
)

func TestOrderClient_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_http.NewMockHttpClient(ctrl)

	type want struct {
		orders models.ProductionOrderPage
		err    error
	}
	type clientCall struct {
		orderApiUrl string
		times       int
		response    *http.Response
		err         error
	}
	tests := []struct {
		name string
		want
		clientCall
	}{
		{
			name: "should fail to get orders when http client returns error",
			want: want{
				orders: models.ProductionOrderPage{},
				err:    errors.New("internal error"),
			},
			clientCall: clientCall{
				orderApiUrl: "/orders",
				times:       1,
				response:    &http.Response{},
				err:         errors.New("internal error"),
			},
		},
		{
			name: "should fail to get orders when response is non-2xx",
			want: want{
				orders: models.ProductionOrderPage{},
				err:    errors.New("resp status code non-2xx"),
			},
			clientCall: clientCall{
				orderApiUrl: "/orders",
				times:       1,
				response: &http.Response{
					StatusCode: 500,
				},
				err: nil,
			},
		},
		{
			name: "should fail to get orders, when resp body is invalid",
			want: want{
				orders: models.ProductionOrderPage{},
				err:    errors.New("failed reading"),
			},
			clientCall: clientCall{
				orderApiUrl: "/orders",
				times:       1,
				response: &http.Response{
					StatusCode: 200,
					Body:       &BrokenReader{},
				},
				err: nil,
			},
		},
		{
			name: "should get orders successfully",
			want: want{
				orders: models.ProductionOrderPage{
					Results: []models.ProductionOrder{createOrder()},
					Next:    new(int),
				},
				err: nil,
			},
			clientCall: clientCall{
				orderApiUrl: "/orders",
				times:       1,
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"results":[{"id":123,"items":[{"id":999,"quantity":1,"type":"UNIT","product":{"id":222,"name":"Batata Frita","skuId":"333","description":"Batata canoa","category":"Acompanhamento","price":9.99,"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}}],"coupon":"APP10","totalAmount":9.99,"status":"PAID","createdAt":"0001-01-01T00:00:00Z","customerCPF":"111222333444"}],"next":0}`)),
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		httpClient.EXPECT().DoGet(gomock.Eq(tt.clientCall.orderApiUrl)).
			Times(tt.clientCall.times).
			Return(tt.clientCall.response, tt.clientCall.err)

		orderClient := NewOrderClient(httpClient, "/orders")
		orders, err := orderClient.GetOrders()

		assert.Equal(t, tt.want.orders, orders)
		assert.Equal(t, tt.want.err, err)
	}
}

func TestOrderClient_UpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_http.NewMockHttpClient(ctrl)

	type args struct {
		orderId string
		status  string
	}
	type want struct {
		err error
	}
	type clientCall struct {
		orderApiUrl string
		times       int
		response    *http.Response
		err         error
	}
	tests := []struct {
		name string
		args
		want
		clientCall
	}{
		{
			name: "should fail to get orders when http client returns error",
			args: args{
				orderId: "123",
				status:  "PAID",
			},
			want: want{
				err: errors.New("internal error"),
			},
			clientCall: clientCall{
				orderApiUrl: "/orders",
				times:       1,
				err:         errors.New("internal error"),
			},
		},
		{
			name: "should fail to get orders when response is non-2xx",
			args: args{
				orderId: "123",
				status:  "PAID",
			},
			want: want{
				err: errors.New("resp status code non-2xx"),
			},
			clientCall: clientCall{
				orderApiUrl: "/orders",
				times:       1,
				response: &http.Response{
					StatusCode: 500,
				},
				err: nil,
			},
		},
		{
			name: "should update order status successfully",
			args: args{
				orderId: "123",
				status:  "PAID",
			},
			want: want{
				err: nil,
			},
			clientCall: clientCall{
				orderApiUrl: "/orders",
				times:       1,
				response: &http.Response{
					StatusCode: 200,
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		httpClient.EXPECT().DoPut(gomock.Any(), gomock.Any()).
			Times(tt.clientCall.times).
			Return(tt.clientCall.response, tt.clientCall.err)

		orderClient := NewOrderClient(httpClient, "/orders")
		err := orderClient.UpdateOrderStatus(tt.args.orderId, tt.args.status)

		assert.Equal(t, tt.want.err, err)
	}
}

type BrokenReader struct{}

func (br *BrokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed reading")
}

func (br *BrokenReader) Close() error {
	return fmt.Errorf("failed closing")
}

func createOrder() models.ProductionOrder {
	return models.ProductionOrder{
		ID:     123,
		Status: "PAID",
		Products: []models.Product{
			{
				Name:        "Batata Frita",
				Description: "Batata canoa",
				Category:    "Acompanhamento",
			},
		},
	}
}

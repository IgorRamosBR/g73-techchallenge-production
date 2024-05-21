package controllers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-production/internal/controllers"
	"github.com/IgorRamosBR/g73-techchallenge-production/internal/core/models"
	mock_usecases "github.com/IgorRamosBR/g73-techchallenge-production/internal/core/usecases/mocks"
	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestOrderController(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "OrderController Suite")
}

var _ = ginkgo.Describe("OrderController", func() {
	var (
		mockCtrl         *gomock.Controller
		mockOrderUseCase *mock_usecases.MockOrderUseCase
		orderController  controllers.OrderController
		router           *gin.Engine
	)

	ginkgo.BeforeEach(func() {
		mockCtrl = gomock.NewController(ginkgo.GinkgoT())
		mockOrderUseCase = mock_usecases.NewMockOrderUseCase(mockCtrl)
		orderController = controllers.NewOrderController(mockOrderUseCase)

		router = gin.Default()
		router.GET("/orders", orderController.GetOrdersHandler)
		router.PUT("/orders/:id/status", orderController.UpdateOrderStatusHandler)
	})

	ginkgo.AfterEach(func() {
		mockCtrl.Finish()
	})

	ginkgo.Describe("GET /orders", func() {
		ginkgo.Context("when there are orders", func() {
			ginkgo.It("should return a list of orders", func() {
				expectedOrders := createOrder()
				mockOrderUseCase.EXPECT().GetOrders().Return(expectedOrders, nil)

				req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				gomega.Expect(w.Code).To(gomega.Equal(http.StatusOK))
			})
		})

		ginkgo.Context("when there is an error fetching orders", func() {
			ginkgo.It("should return an internal server error", func() {
				mockOrderUseCase.EXPECT().GetOrders().Return(models.ProductionOrderPage{}, errors.New("some error"))

				req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				gomega.Expect(w.Code).To(gomega.Equal(http.StatusInternalServerError))
			})
		})
	})

	ginkgo.Describe("PUT /orders/:id/status", func() {
		ginkgo.Context("when the request is valid", func() {
			ginkgo.It("should update the order status", func() {
				orderID := "1"
				mockOrderUseCase.EXPECT().UpdateOrderStatus(orderID, "PAID").Return(nil)

				body := bytes.NewBufferString(`{"status": "PAID"}`)
				req, _ := http.NewRequest(http.MethodPut, "/orders/1/status", body)
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				gomega.Expect(w.Code).To(gomega.Equal(http.StatusNoContent))
			})
		})

		ginkgo.Context("when the request body is invalid", func() {
			ginkgo.It("should return a bad request error", func() {
				body := bytes.NewBufferString(`{"status":`)
				req, _ := http.NewRequest(http.MethodPut, "/orders/1/status", body)
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				gomega.Expect(w.Code).To(gomega.Equal(http.StatusBadRequest))
			})
		})

		ginkgo.Context("when there is an error updating the order status", func() {
			ginkgo.It("should return a bad request error", func() {
				orderID := "1"
				mockOrderUseCase.EXPECT().UpdateOrderStatus(orderID, "PAID").Return(errors.New("some error"))

				body := bytes.NewBufferString(`{"status": "PAID"}`)
				req, _ := http.NewRequest(http.MethodPut, "/orders/1/status", body)
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				gomega.Expect(w.Code).To(gomega.Equal(http.StatusBadRequest))
			})
		})
	})
})

func createOrder() models.ProductionOrderPage {
	return models.ProductionOrderPage{
		Results: []models.ProductionOrder{
			{
				ID:     123,
				Status: "PAID",
				Products: []models.Product{
					{
						Name:        "Batata Frita",
						Description: "Batata canoa",
						Category:    "Acompanhamento",
					},
				},
			},
		},
		Next: new(int),
	}
}

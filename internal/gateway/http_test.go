package gateway_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payment_gateway/internal/gateway"
	"payment_gateway/internal/gateway/mocks"
	"payment_gateway/internal/payment"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

//go:generate mockgen -destination=mocks/mock_gateway.go -package=mocks payment_gateway/internal/gateway PaymentGateway

type HTTPHandlerTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockGateway *mocks.MockPaymentGateway
	handler     *gateway.HTTPHandler
}

func (s *HTTPHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockGateway = mocks.NewMockPaymentGateway(s.ctrl)
	s.handler = gateway.NewHTTPHandler(s.mockGateway)
}

func (s *HTTPHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *HTTPHandlerTestSuite) TestProcessPayment() {
	tests := []struct {
		name           string
		request        gateway.ProcessPaymentRequest
		setupMocks     func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful payment",
			request: gateway.ProcessPaymentRequest{
				MerchantId: "merchant1",
				CardDetails: payment.CardDetails{
					Number:      "1111111111111111",
					ExpiryMonth: 11,
					ExpiryYear:  2026,
					CVV:         "123",
				},
				Amount: payment.Amount{
					Value:    100,
					Currency: "GBP",
				},
			},
			setupMocks: func() {
				s.mockGateway.EXPECT().
					ProcessPayment(gomock.Any(), gomock.Any()).
					Return(&payment.Payment{
						PaymentId: "payment1",
						Status:    payment.PaymentStatusSuccess,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "invalid request",
			request: gateway.ProcessPaymentRequest{},
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			body, err := json.Marshal(tt.request)
			s.Require().NoError(err)

			req := httptest.NewRequest("POST", "/payments", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			s.handler.ProcessPayment(rr, req)

			s.Equal(tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response payment.Payment
				s.Require().NoError(json.NewDecoder(rr.Body).Decode(&response))
				s.NotEmpty(response.PaymentId)
			}
		})
	}
}

func (s *HTTPHandlerTestSuite) TestGetPayment() {
	tests := []struct {
		name           string
		paymentID      string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:      "payment found",
			paymentID: "payment1",
			setupMocks: func() {
				s.mockGateway.EXPECT().
					GetPayment(gomock.Any(), "payment1").
					Return(&payment.Payment{
						PaymentId: "payment1",
						Status:    payment.PaymentStatusSuccess,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "payment not found",
			paymentID: "null",
			setupMocks: func() {
				s.mockGateway.EXPECT().
					GetPayment(gomock.Any(), "null").
					Return(nil, payment.ErrPaymentNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMocks()

			req := httptest.NewRequest("GET", "/payments/"+tt.paymentID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.paymentID})
			rr := httptest.NewRecorder()

			s.handler.GetPayment(rr, req)

			s.Equal(tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response payment.Payment
				s.Require().NoError(json.NewDecoder(rr.Body).Decode(&response))
				s.Equal(tt.paymentID, response.PaymentId)
			}
		})
	}
}

func TestHTTPHandlerSuite(t *testing.T) {
	suite.Run(t, new(HTTPHandlerTestSuite))
}

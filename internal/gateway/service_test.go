package gateway_test

import (
	"context"
	"payment_gateway/internal/bank"
	"payment_gateway/internal/gateway"
	"payment_gateway/internal/gateway/mocks"
	"payment_gateway/internal/payment"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

//go:generate mockgen -destination=mocks/mock_store.go -package=mocks payment_gateway/internal/store PaymentStore
//go:generate mockgen -destination=mocks/mock_bank.go -package=mocks payment_gateway/internal/bank BankService

type ServiceTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	mockStore *mocks.MockPaymentStore
	mockBank  *mocks.MockBankService
	service   gateway.PaymentGateway
}

func (s *ServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockStore = mocks.NewMockPaymentStore(s.ctrl)
	s.mockBank = mocks.NewMockBankService(s.ctrl)
	s.service = gateway.NewService(s.mockStore, s.mockBank)
}

func (s *ServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ServiceTestSuite) createTestRequest() gateway.ProcessPaymentRequest {
	return gateway.ProcessPaymentRequest{
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
	}
}

func (s *ServiceTestSuite) TestProcessPayment() {
	tests := []struct {
		name        string
		request     gateway.ProcessPaymentRequest
		setupMocks  func()
		expectError bool
	}{
		{
			name:    "successful payment",
			request: s.createTestRequest(),
			setupMocks: func() {
				s.mockBank.EXPECT().
					ProcessPayment(gomock.Any(), gomock.Any()).
					Return(&bank.BankResponse{
						TransactionId: "tx1",
						Status:        payment.PaymentStatusSuccess,
					}, nil)

				s.mockStore.EXPECT().
					Save(gomock.Any(), mock.MatchedBy(func(p *payment.Payment) bool {
						return p.CardDetails.MaskedNumber == "************1111" &&
							p.CardDetails.Number == "" &&
							p.CardDetails.CVV == ""
					})).
					Return(nil)
			},
			expectError: false,
		},
		{
			name:    "processing failure",
			request: s.createTestRequest(),
			setupMocks: func() {
				s.mockBank.EXPECT().
					ProcessPayment(gomock.Any(), gomock.Any()).
					Return(nil, payment.ErrProcessingFailed)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMocks()

			result, err := s.service.ProcessPayment(context.Background(), tt.request)

			if tt.expectError {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(payment.PaymentStatusSuccess, result.Status)
				s.Empty(result.CardDetails.Number)
				s.Empty(result.CardDetails.CVV)
			}
		})
	}
}

func (s *ServiceTestSuite) TestGetPayment() {
	tests := []struct {
		name        string
		paymentID   string
		setupMocks  func()
		expectError bool
	}{
		{
			name:      "payment found",
			paymentID: "payment1",
			setupMocks: func() {
				s.mockStore.EXPECT().
					Get(gomock.Any(), "payment1").
					Return(&payment.Payment{
						PaymentId: "payment1",
						Status:    payment.PaymentStatusSuccess,
						CardDetails: payment.CardDetails{
							MaskedNumber: "************1111",
							ExpiryMonth:  11,
							ExpiryYear:   2026,
						},
					}, nil)
			},
			expectError: false,
		},
		{
			name:      "payment not found",
			paymentID: "null",
			setupMocks: func() {
				s.mockStore.EXPECT().
					Get(gomock.Any(), "null").
					Return(nil, payment.ErrPaymentNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMocks()

			result, err := s.service.GetPayment(context.Background(), tt.paymentID)

			if tt.expectError {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(tt.paymentID, result.PaymentId)
				s.NotEmpty(result.CardDetails.MaskedNumber)
				s.Empty(result.CardDetails.Number)
				s.Empty(result.CardDetails.CVV)
			}
		})
	}
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

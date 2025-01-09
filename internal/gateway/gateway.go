package gateway

import (
	"context"
	"payment_gateway/internal/payment"
)

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*payment.Payment, error)
	GetPayment(ctx context.Context, id string) (*payment.Payment, error)
}

type ProcessPaymentRequest struct {
	MerchantId  string              `json:"merchant_id" validate:"required"`
	CardDetails payment.CardDetails `json:"card_details" validate:"required"`
	Amount      payment.Amount      `json:"amount" validate:"required"`
}

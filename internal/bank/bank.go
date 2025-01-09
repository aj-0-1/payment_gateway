package bank

import (
	"context"
	"payment_gateway/internal/payment"
)

type BankService interface {
	ProcessPayment(ctx context.Context, payment *payment.Payment) (*BankResponse, error)
}

type BankResponse struct {
	TransactionId string                `json:"transaction_id"`
	Status        payment.PaymentStatus `json:"status"`
}

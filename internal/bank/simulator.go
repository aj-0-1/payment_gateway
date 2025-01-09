package bank

import (
	"context"
	"payment_gateway/internal/payment"
	"time"

	"github.com/google/uuid"
)

type Simulator struct{}

func NewSimulator() *Simulator {
	return &Simulator{}
}

func (s *Simulator) ProcessPayment(ctx context.Context, p *payment.Payment) (*BankResponse, error) {
	// Simulate a typical payment validation by a bank
	time.Sleep(100 * time.Millisecond)

	return &BankResponse{
		TransactionId: uuid.New().String(),
		Status:        payment.PaymentStatusSuccess,
	}, nil
}

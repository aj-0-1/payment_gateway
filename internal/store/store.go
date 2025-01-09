package store

import (
	"context"
	"payment_gateway/internal/payment"
)

type PaymentStore interface {
	Save(ctx context.Context, payment *payment.Payment) error
	Get(ctx context.Context, id string) (*payment.Payment, error)
}

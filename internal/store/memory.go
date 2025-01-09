package store

import (
	"context"
	"payment_gateway/internal/payment"
	"sync"
)

type MemoryStore struct {
	payments map[string]*payment.Payment
	mu       sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		payments: make(map[string]*payment.Payment),
	}
}

func (ms *MemoryStore) Save(ctx context.Context, p *payment.Payment) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.payments[p.PaymentId] = p
	return nil
}

func (ms *MemoryStore) Get(ctx context.Context, id string) (*payment.Payment, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if payment, ok := ms.payments[id]; ok {
		return payment, nil
	}
	return nil, payment.ErrPaymentNotFound
}

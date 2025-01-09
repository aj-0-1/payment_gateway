package gateway

import (
	"context"
	"payment_gateway/internal/bank"
	"payment_gateway/internal/payment"
	"payment_gateway/internal/store"
	"time"

	"github.com/google/uuid"
)

type service struct {
	store store.PaymentStore
	bank  bank.BankService
}

func NewService(store store.PaymentStore, bank bank.BankService) PaymentGateway {
	return &service{
		store: store,
		bank:  bank,
	}
}

func (s *service) ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*payment.Payment, error) {
	p := &payment.Payment{
		PaymentId:   uuid.New().String(),
		MerchantId:  req.MerchantId,
		CardDetails: req.CardDetails,
		Amount:      req.Amount,
		Status:      payment.PaymentStatusPending,
		CreatedAt:   time.Now(),
	}

	bankResp, err := s.bank.ProcessPayment(ctx, p)
	if err != nil {
		return nil, err
	}

	p.Status = bankResp.Status
	p.CardDetails.MaskedNumber = "************" + p.CardDetails.Number[12:]
	p.CardDetails.Number = ""
	p.CardDetails.CVV = ""

	if err := s.store.Save(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) GetPayment(ctx context.Context, id string) (*payment.Payment, error) {
	return s.store.Get(ctx, id)
}

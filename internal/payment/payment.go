package payment

import (
	"time"
)

type Payment struct {
	PaymentId   string        `json:"payment_id"`
	MerchantId  string        `json:"merchant_id"`
	CardDetails CardDetails   `json:"card_details,omitempty"`
	Amount      Amount        `json:"amount"`
	Status      PaymentStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
}

type CardDetails struct {
	Number       string `json:"number" validate:"required,min=16,max=16"`
	MaskedNumber string `json:"masked_number,omitempty"`
	ExpiryMonth  int    `json:"expiry_month" validate:"required,min=1,max=12"`
	ExpiryYear   int    `json:"expiry_year" validate:"required,min=2025"`
	CVV          string `json:"cvv" validate:"required,min=3,max=4"`
}

type Amount struct {
	Value    int64  `json:"value" validate:"required,gt=0"`
	Currency string `json:"currency" validate:"required,len=3"`
}

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

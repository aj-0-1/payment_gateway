package payment

import "errors"

var (
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrProcessingFailed = errors.New("payment processing failed")
)

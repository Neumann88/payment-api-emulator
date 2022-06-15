package contracts

import (
	"context"

	"github.com/Neumann88/payment-api-emulator/internal/entity"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, input entity.PaymentInput) (int64, error)
	UpdateStatus(ctx context.Context, input entity.PaymentStatus) (int64, error)
	GetStatus(ctx context.Context, paymentID int64) (string, error)
	GetPayments(ctx context.Context, input entity.PaymentUser) ([]entity.Payment, error)
	CancelPayment(ctx context.Context, paymentID int64) (int64, error)
}

type PaymentUseCase interface {
	CreatePayment(ctx context.Context, input entity.PaymentInput) (int64, error)
	UpdateStatus(ctx context.Context, input entity.PaymentStatus) error
	GetStatus(ctx context.Context, paymentID int64) (string, error)
	GetPayments(ctx context.Context, input entity.PaymentUser) ([]entity.Payment, error)
	CancelPayment(ctx context.Context, paymentID int64) error
}

package payment

import (
	"context"

	"github.com/gorilla/mux"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, input PaymentInput) (int64, error)
	UpdateStatus(ctx context.Context, input PaymentStatus) (int64, error)
	GetStatus(ctx context.Context, paymentID int64) (string, error)
	GetPayments(ctx context.Context, input PaymentUser) ([]PaymentData, error)
}

type PaymentUsecase interface {
	CreatePayment(ctx context.Context, input PaymentInput) (int64, error)
	UpdateStatus(ctx context.Context, input PaymentStatus) error
	GetStatus(ctx context.Context, paymentID int64) (string, error)
	GetPayments(ctx context.Context, input PaymentUser) ([]PaymentData, error)
}

type PaymentController interface {
	Register(router *mux.Router)
}

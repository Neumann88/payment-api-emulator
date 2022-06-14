package payment

import (
	"context"
)

type paymentRepository interface {
	createPayment(ctx context.Context, input paymentInput) (int64, error)
	updateStatus(ctx context.Context, input paymentStatus) (int64, error)
	getStatus(ctx context.Context, paymentID int64) (string, error)
	getPayments(ctx context.Context, input paymentUser) ([]payment, error)
	cancelPayment(ctx context.Context, paymentID int64) (int64, error)
}

type paymentUseCase interface {
	createPayment(ctx context.Context, input paymentInput) (int64, error)
	updateStatus(ctx context.Context, input paymentStatus) error
	getStatus(ctx context.Context, paymentID int64) (string, error)
	getPayments(ctx context.Context, input paymentUser) ([]payment, error)
	cancelPayment(ctx context.Context, paymentID int64) error
}

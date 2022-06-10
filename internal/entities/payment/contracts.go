package payment

import (
	"context"

	"github.com/gorilla/mux"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, input PaymentInput) (int64, error)
}

type PaymentUsecase interface {
	CreatePayment(ctx context.Context, input PaymentInput) (int64, error)
}

type PaymentController interface {
	Register(router *mux.Router)
}

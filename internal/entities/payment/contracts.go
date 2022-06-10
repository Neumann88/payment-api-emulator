package payment

import (
	"context"

	"github.com/gorilla/mux"
)

type PaymentRepository interface {
	createPayment(ctx context.Context, input Payment) (int64, error)
}

type PaymentUsecase interface {
	createPayment(ctx context.Context, input Payment) (int64, error)
}

type PaymentController interface {
	Register(router *mux.Router)
}

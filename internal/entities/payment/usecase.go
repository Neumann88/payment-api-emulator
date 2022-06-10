package payment

import "context"

type Usecase struct {
	repo PaymentRepository
}

func NewPaymentUsecase(r PaymentRepository) *Usecase {
	return &Usecase{
		repo: r,
	}
}

func (a *Usecase) createPayment(ctx context.Context, input Payment) (int64, error) {
	return a.repo.createPayment(ctx, input)
}

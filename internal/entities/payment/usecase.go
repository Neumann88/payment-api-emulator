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

func (a *Usecase) CreatePayment(ctx context.Context, input PaymentInput) (int64, error) {
	return a.repo.CreatePayment(ctx, input)
}

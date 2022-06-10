package usecase

import (
	"github.com/Neumann88/payment-api-emulator/internal/entities/payment"
	"github.com/Neumann88/payment-api-emulator/internal/repository"
)

type Usecase struct {
	PaymentUsecase payment.PaymentUsecase
}

func NewUsecase(repos *repository.Repository) *Usecase {
	return &Usecase{
		PaymentUsecase: payment.NewPaymentUsecase(repos.PaymentRepository),
	}
}

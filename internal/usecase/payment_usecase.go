package usecase

import (
	"context"
	"errors"
	"sync"

	"github.com/Neumann88/payment-api-emulator/internal/contracts"
	"github.com/Neumann88/payment-api-emulator/internal/entity"
)

type PaymentUseCase struct {
	repo contracts.PaymentRepository
}

func NewPaymentUseCase(repo contracts.PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{
		repo: repo,
	}
}

func (u *PaymentUseCase) CreatePayment(ctx context.Context, input entity.PaymentInput) (int64, error) {
	paymentID, err := u.repo.CreatePayment(
		ctx,
		input,
	)

	if err != nil {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			err = u.UpdateStatus(
				ctx,
				entity.PaymentStatus{
					ID:     paymentID,
					Status: entity.StatusError,
				},
			)

			wg.Done()
		}()
		wg.Wait()

		return 0, err
	}

	return paymentID, nil
}

func (u *PaymentUseCase) UpdateStatus(ctx context.Context, input entity.PaymentStatus) error {
	row, err := u.repo.UpdateStatus(
		ctx,
		input,
	)

	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("payment-usecase-updateStatus, terminal status")
	}

	return nil
}

func (u *PaymentUseCase) GetStatus(ctx context.Context, paymentID int64) (string, error) {
	return u.repo.GetStatus(
		ctx,
		paymentID,
	)
}

func (u *PaymentUseCase) GetPayments(ctx context.Context, input entity.PaymentUser) ([]entity.Payment, error) {
	return u.repo.GetPayments(
		ctx,
		input,
	)
}

func (u *PaymentUseCase) CancelPayment(ctx context.Context, paymentID int64) error {
	row, err := u.repo.CancelPayment(
		ctx,
		paymentID,
	)

	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("payment-usecase-cancelPayment, terminal status")
	}

	return nil
}

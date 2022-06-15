package usecase

import (
	"context"
	"errors"
	"sync"

	"github.com/Neumann88/payment-api-emulator/internal/contracts"
	"github.com/Neumann88/payment-api-emulator/internal/entity"
)

type paymentUseCase struct {
	repo contracts.PaymentRepository
}

func NewPaymentUseCase(repo contracts.PaymentRepository) *paymentUseCase {
	return &paymentUseCase{
		repo: repo,
	}
}

func (u *paymentUseCase) CreatePayment(ctx context.Context, input entity.PaymentInput) (int64, error) {
	paymentID, err := u.repo.CreatePayment(
		ctx,
		input,
	)

	if err != nil {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			_ = u.UpdateStatus(
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

func (u *paymentUseCase) UpdateStatus(ctx context.Context, input entity.PaymentStatus) error {
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

func (u *paymentUseCase) GetStatus(ctx context.Context, paymentID int64) (string, error) {
	return u.repo.GetStatus(
		ctx,
		paymentID,
	)
}

func (u *paymentUseCase) GetPayments(ctx context.Context, input entity.PaymentUser) ([]entity.Payment, error) {
	return u.repo.GetPayments(
		ctx,
		input,
	)
}

func (u *paymentUseCase) CancelPayment(ctx context.Context, paymentID int64) error {
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

package payment

import (
	"context"
	"fmt"
	"sync"
)

type useCase struct {
	repo paymentRepository
}

func NewPaymentUseCase(repo paymentRepository) *useCase {
	return &useCase{
		repo: repo,
	}
}

func (u *useCase) createPayment(ctx context.Context, input paymentInput) (int64, error) {
	paymentID, err := u.repo.createPayment(
		ctx,
		input,
	)

	if err != nil {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			_ = u.updateStatus(
				ctx,
				paymentStatus{
					ID:     paymentID,
					Status: statusError,
				},
			)

			wg.Done()
		}()
		wg.Wait()

		return 0, err
	}

	return paymentID, nil
}

func (u *useCase) updateStatus(ctx context.Context, input paymentStatus) error {
	status, err := u.repo.getStatus(
		ctx,
		input.ID,
	)

	if err != nil {
		return fmt.Errorf("payment-usecase-updateStatus %s", err.Error())
	}

	if status == statusSuccess || status == statusFailure {
		return fmt.Errorf("payment-usecase-updateStatus, terminal status: %s", status)
	}

	r, err := u.repo.updateStatus(
		ctx,
		input,
	)

	if err != nil {
		return err
	}

	err = checkTerminalStatusRow(r)

	if err != nil {
		return fmt.Errorf("payment-usecase-updateStatus, %s", err.Error())
	}

	return nil
}

func (u *useCase) getStatus(ctx context.Context, paymentID int64) (string, error) {
	return u.repo.getStatus(
		ctx,
		paymentID,
	)
}

func (u *useCase) getPayments(ctx context.Context, input paymentUser) ([]payment, error) {
	return u.repo.getPayments(
		ctx,
		input,
	)
}

func (u *useCase) deletePayment(ctx context.Context, paymentID int64) error {
	r, err := u.repo.deletePayment(
		ctx,
		paymentID,
	)

	if err != nil {
		return err
	}

	err = checkTerminalStatusRow(r)

	if err != nil {
		return fmt.Errorf("Payment-Usecase-deletePayment, %s", err.Error())
	}

	return nil
}

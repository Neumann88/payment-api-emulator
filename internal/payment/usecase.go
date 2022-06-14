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

// TODO: дубликат
func (u *useCase) updateStatus(ctx context.Context, input paymentStatus) error {
	errorExeption := make(chan error)
	uCtx, cancel := context.WithCancel(ctx)

	go func() {
		status, err := u.repo.getStatus(
			uCtx,
			input.ID,
		)

		if err != nil {
			return
		}

		if status == statusSuccess || status == statusFailure {
			errorExeption <- fmt.Errorf("payment-usecase-updateStatus-getStatus, terminal status %s", status)
		}
	}()

	go func() {
		r, err := u.repo.updateStatus(
			uCtx,
			input,
		)

		if err != nil {
			errorExeption <- err
			return
		}

		err = checkTerminalStatusRow(r)

		if err != nil {
			errorExeption <- fmt.Errorf("payment-usecase-updateStatus, %s", err.Error())
		}

		cancel()
	}()

	select {
	case <-uCtx.Done():
		return nil
	case err := <-errorExeption:
		return err
	}
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

// TODO: дубликат
func (u *useCase) deletePayment(ctx context.Context, paymentID int64) error {
	errorExeption := make(chan error)
	dCtx, cancel := context.WithCancel(ctx)

	go func() {
		status, err := u.repo.getStatus(
			dCtx,
			paymentID,
		)

		if err != nil {
			return
		}

		if status == statusSuccess || status == statusFailure {
			errorExeption <- fmt.Errorf("payment-usecase-deletePayment-getStatus, terminal status %s", status)
		}
	}()

	go func() {
		r, err := u.repo.deletePayment(
			dCtx,
			paymentID,
		)

		if err != nil {
			errorExeption <- err
		}

		err = checkTerminalStatusRow(r)

		if err != nil {
			errorExeption <- fmt.Errorf("payment-usecase-deletePayment, %s", err.Error())
		}

		cancel()
	}()

	select {
	case <-dCtx.Done():
		return nil
	case err := <-errorExeption:
		return err
	}
}

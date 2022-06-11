package payment

import (
	"context"
	"fmt"
	"sync"
)

type Usecase struct {
	repo PaymentRepository
}

func NewPaymentUsecase(r PaymentRepository) *Usecase {
	return &Usecase{
		repo: r,
	}
}

func (u *Usecase) CreatePayment(ctx context.Context, input PaymentInput) (int64, error) {
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
				PaymentStatus{
					ID:     paymentID,
					Status: StatusError,
				},
			)

			wg.Done()
		}()
		wg.Wait()

		return 0, err
	}

	return paymentID, nil
}

func (u *Usecase) UpdateStatus(ctx context.Context, input PaymentStatus) error {
	status, err := u.repo.GetStatus(
		ctx,
		input.ID,
	)

	if err != nil {
		return fmt.Errorf("Payment-Usecase-UpdateStatus %s", err.Error())
	}

	if status == StatusSuccess || status == StatusFailure {
		return fmt.Errorf("Payment-Usecase-UpdateStatus, terminal status: %s", status)
	}

	if input.Status == StatusError {
		return fmt.Errorf("Payment-Usecase-UpdateStatus, invalid status: %s", input.Status)
	}

	r, err := u.repo.UpdateStatus(
		ctx,
		input,
	)

	if err != nil {
		return err
	}

	err = checkTerminalStatusRow(r)

	if err != nil {
		return fmt.Errorf("Payment-Usecase-UpdateStatus, %s", err.Error())
	}

	return nil
}

func (u *Usecase) GetStatus(ctx context.Context, paymentID int64) (string, error) {
	return u.repo.GetStatus(
		ctx,
		paymentID,
	)
}

func (u *Usecase) GetPayments(ctx context.Context, input PaymentUser) ([]Payment, error) {
	return u.repo.GetPayments(
		ctx,
		input,
	)
}

func (u *Usecase) CancelPayment(ctx context.Context, paymentID int64) error {
	r, err := u.repo.CancelPayment(
		ctx,
		paymentID,
	)

	if err != nil {
		return err
	}

	err = checkTerminalStatusRow(r)

	if err != nil {
		return fmt.Errorf("Payment-Usecase-CancelPayment, %s", err.Error())
	}

	return nil
}

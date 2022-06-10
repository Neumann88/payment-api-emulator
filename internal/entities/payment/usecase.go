package payment

import (
	"context"
	"errors"
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
	paymentID, err := u.repo.CreatePayment(ctx, input)

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
	status, err := u.repo.GetStatus(ctx, input.ID)

	if err != nil {
		return fmt.Errorf("Payment-Usecase-Repository-GetStatus %s", err.Error())
	}

	if status == StatusSuccess || status == StatusFailure {
		return fmt.Errorf("Payment-Usecase-Repository-GetStatus, terminal status: %s", status)
	}

	if input.Status == StatusError {
		return fmt.Errorf("Payment-Usecase-UpdateStatus, invalid status: %s", input.Status)
	}

	r, err := u.repo.UpdateStatus(ctx, input)

	if err != nil {
		return err
	}

	if r == 0 {
		return errors.New("Payment-Usecase-UpdateStatus, zero rows")
	}

	return nil
}

func (u *Usecase) GetStatus(ctx context.Context, paymentID int64) (string, error) {
	return u.repo.GetStatus(ctx, paymentID)
}

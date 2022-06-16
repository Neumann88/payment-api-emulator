package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Neumann88/payment-api-emulator/internal/entity"
)

const payments = "payments"

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, input entity.PaymentInput) (int64, error) {
	const format = `INSERT INTO %s (user_id, user_email, amount, currency)
						VALUES ($1, $2, $3, $4)
					RETURNING id`

	query := fmt.Sprintf(
		format,
		payments,
	)

	row := r.db.QueryRowContext(
		ctx,
		query,
		input.UserID,
		input.UserEmail,
		input.Amount,
		input.Currency,
	)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("payment-repository-createPayment, %w", err)
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("payment-repository-createPayment, no result")
		}

		return 0, fmt.Errorf("payment-repository-createPayment, %w", err)
	}

	return id, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, input entity.PaymentStatus) (int64, error) {
	const format = `UPDATE %s SET status = $1
						WHERE id = $2
						AND status NOT IN ($3, $4)`

	query := fmt.Sprintf(
		format,
		payments,
	)

	rows, err := r.db.ExecContext(
		ctx,
		query,
		input.Status,
		input.ID,
		entity.StatusSuccess,
		entity.StatusFailure,
	)

	if err != nil {
		return 0, fmt.Errorf("payment-reposiroty-updateStatus, %w", err)
	}

	return rows.RowsAffected()
}

func (r *PaymentRepository) GetStatus(ctx context.Context, paymentID int64) (string, error) {
	const format = `SELECT status from %s
						WHERE id = $1`

	query := fmt.Sprintf(
		format,
		payments,
	)

	rows := r.db.QueryRowContext(
		ctx,
		query,
		paymentID,
	)

	var status string
	if err := rows.Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("payment-repository-createPayment, no result")
		}

		return "", fmt.Errorf("payment-reposiroty-getStatus, %w", err)
	}

	return status, nil
}

func (r *PaymentRepository) GetPayments(ctx context.Context, input entity.PaymentUser) ([]entity.Payment, error) {
	var arg string
	var value interface{}

	if input.UserEmail != "" && input.UserID == 0 {
		arg = "user_email"
		value = input.UserEmail
	}

	if input.UserID != 0 && input.UserEmail == "" {
		arg = "user_id"
		value = input.UserID
	}

	const format = `SELECT
						id,
						user_id,
						user_email,
						currency,
						amount,
						created_at,
						updated_at,
						status
					from %s
						WHERE %s = $1`

	query := fmt.Sprintf(
		format,
		payments,
		arg,
	)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		value,
	)

	if err != nil {
		return []entity.Payment{}, fmt.Errorf("payment-reposiroty-getPayments, %w", err)
	}

	defer rows.Close()

	output := make([]entity.Payment, 0)
	for rows.Next() {
		value := entity.Payment{}

		err := rows.Scan(
			&value.ID,
			&value.UserID,
			&value.UserEmail,
			&value.Currency,
			&value.Amount,
			&value.CreatedAt,
			&value.UpdatedAt,
			&value.Status,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return []entity.Payment{}, errors.New("payment-repository-createPayment, no result")
			}

			return []entity.Payment{}, fmt.Errorf("payment-reposiroty-getPayments, %w", err)
		}

		output = append(output, value)
	}

	if err := rows.Err(); err != nil {
		return []entity.Payment{}, fmt.Errorf("payment-reposiroty-getPayments, %w", err)
	}

	return output, nil
}

func (r *PaymentRepository) CancelPayment(ctx context.Context, paymentID int64) (int64, error) {
	const format = `UPDATE %s SET status = $1
						WHERE id = $2
						AND status NOT IN ($3, $4)`

	query := fmt.Sprintf(
		format,
		payments,
	)

	rows, err := r.db.ExecContext(
		ctx,
		query,
		entity.StatusCanceled,
		paymentID,
		entity.StatusSuccess,
		entity.StatusFailure,
	)

	if err != nil {
		return 0, fmt.Errorf("payment-repository-cancelPayment, %w", err)
	}

	return rows.RowsAffected()
}

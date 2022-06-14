package payment

import (
	"context"
	"database/sql"
	"fmt"
)

const payments = "payments"

type repository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) createPayment(ctx context.Context, input paymentInput) (int64, error) {
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
		return 0, fmt.Errorf("payment-repository-createPayment, %s", err.Error())
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("payment-repository-createPayment, %s", "no result")
		}

		return 0, fmt.Errorf("payment-repository-createPayment, %s", err.Error())
	}

	return id, nil
}

func (r *repository) updateStatus(ctx context.Context, input paymentStatus) (int64, error) {
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
		statusSuccess,
		statusFailure,
	)

	if err != nil {
		return 0, fmt.Errorf("payment-reposiroty-updateStatus, %s", err.Error())
	}

	return rows.RowsAffected()
}

func (r *repository) getStatus(ctx context.Context, paymentID int64) (string, error) {
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
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("payment-repository-createPayment, %s", "no result")
		}

		return "", fmt.Errorf("payment-reposiroty-getStatus, %s", err.Error())
	}

	return status, nil
}

func (r *repository) getPayments(ctx context.Context, input paymentUser) ([]payment, error) {
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
		return []payment{}, fmt.Errorf("payment-reposiroty-getPayments, %s", err.Error())
	}

	defer rows.Close()

	output := make([]payment, 0)
	for rows.Next() {
		value := payment{}

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
			if err == sql.ErrNoRows {
				return []payment{}, fmt.Errorf("payment-repository-createPayment, %s", "no result")
			}

			return []payment{}, fmt.Errorf("payment-reposiroty-getPayments, %s", err.Error())
		}

		output = append(output, value)
	}

	if err := rows.Err(); err != nil {
		return []payment{}, fmt.Errorf("payment-reposiroty-getPayments, %s", err.Error())
	}

	return output, nil
}

func (r *repository) deletePayment(ctx context.Context, paymentID int64) (int64, error) {
	const format = `DELETE FROM %s
						WHERE id = $1
						AND status NOT IN ($2, $3)`

	query := fmt.Sprintf(
		format,
		payments,
	)

	rows, err := r.db.ExecContext(
		ctx,
		query,
		paymentID,
		statusSuccess,
		statusFailure,
	)

	if err != nil {
		return 0, fmt.Errorf("Payment-Repository-deletePayment, %s", err.Error())
	}

	return rows.RowsAffected()
}

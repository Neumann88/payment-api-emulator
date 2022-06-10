package payment

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreatePayment(ctx context.Context, input PaymentInput) (int64, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, user_email, amount, currency)
			VALUES ($1, $2, $3, $4)
		RETURNING id`,
		Payments,
	)

	row := r.db.QueryRowContext(
		ctx,
		query,
		input.UserID,
		input.UserEmail,
		input.Amount,
		input.Currency,
	)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("Payment-Repository-CreatePayment, %s", err.Error())
	}

	return id, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, input PaymentStatus) (int64, error) {
	query := fmt.Sprintf(
		`UPDATE %s
			SET status = $1
				WHERE id = $2
			AND status NOT IN ($3, $4)`,
		Payments,
	)

	rows, err := r.db.ExecContext(
		ctx,
		query,
		input.Status,
		input.ID,
		StatusSuccess,
		StatusFailure,
	)

	if err != nil {
		return 0, fmt.Errorf("Payment-Reposiroty-UpdateStatus, %s", err.Error())
	}

	return rows.RowsAffected()
}

func (r *Repository) GetStatus(ctx context.Context, paymentID int64) (string, error) {
	query := fmt.Sprintf(
		`SELECT status from %s
			WHERE id = $1`,
		Payments,
	)

	rows := r.db.QueryRowContext(
		ctx,
		query,
		paymentID,
	)

	var status string
	if err := rows.Scan(&status); err != nil {
		return "", fmt.Errorf("Payment-Reposiroty-GetStatusByID, %s", err.Error())
	}

	return status, nil
}

func (r *Repository) GetPayments(ctx context.Context, input PaymentUser) ([]PaymentData, error) {
	var arg string
	var value interface{}

	if input.UserEmail != "" && input.UserID == 0 {
		arg = User_Email
		value = input.UserEmail
	}

	if input.UserID != 0 && input.UserEmail == "" {
		arg = User_ID
		value = input.UserID
	}

	query := fmt.Sprintf(
		`SELECT * from %s
			WHERE %s = $1`,
		Payments,
		arg,
	)

	var output []PaymentData

	err := r.db.SelectContext(
		ctx,
		&output,
		query,
		value,
	)

	if err != nil {
		return []PaymentData{}, fmt.Errorf("Payment-Reposiroty-GetPayments, %s", err.Error())
	}

	return output, nil
}

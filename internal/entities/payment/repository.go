package payment

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

const payments = "payments"

func NewPaymentRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreatePayment(ctx context.Context, input PaymentInput) (int64, error) {
	query := fmt.Sprintf(
		"INSERT INTO %s (user_id, user_email, amount, currency) VALUES ($1, $2, $3, $4) RETURNING id",
		payments,
	)

	var id int64
	row := r.db.QueryRowContext(
		ctx,
		query,
		input.UserID,
		input.UserEmail,
		input.Amount,
		input.Currency,
	)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("Payment-Repository-CreatePayment, %s", err.Error())
	}

	return id, nil
}

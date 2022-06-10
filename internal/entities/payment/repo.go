package payment

import (
	"context"

	"github.com/Neumann88/payment-api-emulator/pkg/loggin"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	logger loggin.ILogger
	db     *sqlx.DB
}

func NewPaymentRepository(l loggin.ILogger, db *sqlx.DB) *Repository {
	return &Repository{
		logger: l,
		db:     db,
	}
}

func (a *Repository) createPayment(ctx context.Context, input Payment) (int64, error) {
	panic("implement me")
}

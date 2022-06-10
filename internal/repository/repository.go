package repository

import (
	"github.com/Neumann88/payment-api-emulator/internal/entities/payment"
	"github.com/Neumann88/payment-api-emulator/pkg/loggin"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	PaymentRepository payment.PaymentRepository
}

func NewRepository(l loggin.ILogger, db *sqlx.DB) *Repository {
	return &Repository{
		PaymentRepository: payment.NewPaymentRepository(l, db),
	}
}

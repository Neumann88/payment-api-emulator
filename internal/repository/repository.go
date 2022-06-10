package repository

import (
	"github.com/Neumann88/payment-api-emulator/internal/entities/payment"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	PaymentRepository payment.PaymentRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		PaymentRepository: payment.NewPaymentRepository(db),
	}
}

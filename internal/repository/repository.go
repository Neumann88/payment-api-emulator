package repository

import (
	"database/sql"

	"github.com/Neumann88/payment-api-emulator/internal/entities/payment"
)

type Repository struct {
	PaymentRepository payment.PaymentRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		PaymentRepository: payment.NewPaymentRepository(db),
	}
}

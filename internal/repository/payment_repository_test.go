package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"

	"github.com/Neumann88/payment-api-emulator/internal/entity"
)

func TestCreatePayment(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  entity.PaymentInput
		expect int64
		err    error
	}{
		{
			name: "Create payment",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				dbMock.ExpectQuery("INSERT INTO payments").
					WithArgs(1, "user_email", 10.0, "currency").
					WillReturnRows(rows)
			},
			input: entity.PaymentInput{
				UserID:    1,
				UserEmail: "user_email",
				Amount:    10.0,
				Currency:  "currency",
			},
			expect: 1,
			err:    nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectQuery("INSERT INTO payments").
					WithArgs(0, "", 0.0, "").
					WillReturnError(errors.New("insert error"))
			},
			input: entity.PaymentInput{
				UserID:    0,
				UserEmail: "",
				Amount:    0.0,
				Currency:  "",
			},
			err: errors.New("insert error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreatePayment(
				context.TODO(),
				entity.PaymentInput{
					UserID:    tt.input.UserID,
					UserEmail: tt.input.UserEmail,
					Amount:    tt.input.Amount,
					Currency:  tt.input.Currency,
				})

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  entity.PaymentStatus
		expect int64
		err    error
	}{
		{
			name: "Update status to success",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("success", 1, entity.StatusSuccess, entity.StatusFailure).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: entity.PaymentStatus{
				ID:     1,
				Status: "success",
			},
			expect: 1,
			err:    nil,
		},
		{
			name: "Update status to failure",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("failure", 1, entity.StatusSuccess, entity.StatusFailure).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: entity.PaymentStatus{
				ID:     1,
				Status: "failure",
			},
			expect: 1,
			err:    nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("failure", 1, entity.StatusSuccess, entity.StatusFailure).
					WillReturnError(errors.New("update error"))
			},
			input: entity.PaymentStatus{
				ID:     1,
				Status: "failure",
			},
			err: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.UpdateStatus(
				context.TODO(),
				entity.PaymentStatus{
					ID:     tt.input.ID,
					Status: tt.input.Status,
				})

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestGetStatus(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  int64
		expect string
		err    error
	}{
		{
			name: "Get payment status",
			mock: func() {
				row := sqlmock.NewRows([]string{"status"}).AddRow("new")

				dbMock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnRows(row)
			},
			input:  1,
			expect: "new",
			err:    nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnError(errors.New("not found"))
			},
			input: 1,
			err:   errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetStatus(
				context.TODO(),
				tt.input,
			)

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestGetPayments(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  entity.PaymentUser
		expect []entity.Payment
		err    error
	}{
		{
			name: "Get user payments by ID",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "user_email", "currency", "amount", "created_at", "updated_at", "status"}).
					AddRow(1, 1, "user_email", "currency", 10.0, "created_at", "updated_at", "status").
					AddRow(2, 2, "user_email", "currency", 10.0, "created_at", "updated_at", "status")

				dbMock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnRows(rows)
			},
			input: entity.PaymentUser{
				UserID: 1,
			},
			expect: []entity.Payment{
				{
					ID:        1,
					UserID:    1,
					Amount:    10.0,
					UserEmail: "user_email",
					Currency:  "currency",
					CreatedAt: "created_at",
					UpdatedAt: "updated_at",
					Status:    "status",
				},
				{
					ID:        2,
					UserID:    2,
					Amount:    10.0,
					UserEmail: "user_email",
					Currency:  "currency",
					CreatedAt: "created_at",
					UpdatedAt: "updated_at",
					Status:    "status",
				},
			},
			err: nil,
		},
		{
			name: "Get user payments by Email",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "user_email", "currency", "amount", "created_at", "updated_at", "status"}).
					AddRow(1, 1, "user_email", "currency", 10.0, "created_at", "updated_at", "status").
					AddRow(2, 2, "user_email", "currency", 10.0, "created_at", "updated_at", "status")

				dbMock.ExpectQuery("SELECT").
					WithArgs("email").
					WillReturnRows(rows)
			},
			input: entity.PaymentUser{
				UserEmail: "email",
			},
			expect: []entity.Payment{
				{
					ID:        1,
					UserID:    1,
					Amount:    10.0,
					UserEmail: "user_email",
					Currency:  "currency",
					CreatedAt: "created_at",
					UpdatedAt: "updated_at",
					Status:    "status",
				},
				{
					ID:        2,
					UserID:    2,
					Amount:    10.0,
					UserEmail: "user_email",
					Currency:  "currency",
					CreatedAt: "created_at",
					UpdatedAt: "updated_at",
					Status:    "status",
				},
			},
			err: nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnError(errors.New("not found"))
			},
			input: entity.PaymentUser{
				UserID: 1,
			},
			err: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetPayments(
				context.TODO(),
				entity.PaymentUser{
					UserID:    tt.input.UserID,
					UserEmail: tt.input.UserEmail,
				},
			)

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestCancelPayment(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  int64
		expect int64
		err    error
	}{
		{
			name: "Cancel payment",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("canceled", 1, entity.StatusSuccess, entity.StatusFailure).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:  1,
			expect: 1,
			err:    nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("canceled", 1, entity.StatusSuccess, entity.StatusFailure).
					WillReturnError(errors.New("update error"))
			},
			input: 1,
			err:   errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CancelPayment(
				context.TODO(),
				tt.input,
			)

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

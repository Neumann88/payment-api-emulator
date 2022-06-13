package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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
		input  paymentInput
		expect int64
		err    error
	}{
		{
			name: "Create payment",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				dbMock.ExpectQuery("INSERT INTO payments").
					WithArgs(1, "user_email", 10.5, "currency").
					WillReturnRows(rows)
			},
			input: paymentInput{
				UserID:    1,
				UserEmail: "user_email",
				Amount:    10.5,
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
			input: paymentInput{
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

			got, err := r.createPayment(
				context.TODO(),
				paymentInput{
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
		input  paymentStatus
		expect int64
		err    error
	}{
		{
			name: "Update status to success",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("success", 1, statusSuccess, statusFailure).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: paymentStatus{
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
					WithArgs("failure", 1, statusSuccess, statusFailure).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: paymentStatus{
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
					WithArgs("failure", 1, statusSuccess, statusFailure).
					WillReturnError(errors.New("update error"))
			},
			input: paymentStatus{
				ID:     1,
				Status: statusFailure,
			},
			err: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.updateStatus(
				context.TODO(),
				paymentStatus{
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

			got, err := r.getStatus(
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
		input  paymentUser
		expect []payment
		err    error
	}{
		{
			name: "Get user payments by ID",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "user_email", "currency", "amount", "created_at", "updated_at", "status"}).
					AddRow(1, 1, "user_email", "currency", 10.5, "created_at", "updated_at", "status").
					AddRow(2, 2, "user_email", "currency", 10.5, "created_at", "updated_at", "status")

				dbMock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnRows(rows)
			},
			input: paymentUser{
				UserID: 1,
			},
			expect: []payment{
				{1, 1, 10.5, "user_email", "currency", "created_at", "updated_at", "status"},
				{2, 2, 10.5, "user_email", "currency", "created_at", "updated_at", "status"},
			},
			err: nil,
		},
		{
			name: "Get user payments by Email",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "user_email", "currency", "amount", "created_at", "updated_at", "status"}).
					AddRow(1, 1, "user_email", "currency", 10.5, "created_at", "updated_at", "status").
					AddRow(2, 2, "user_email", "currency", 10.5, "created_at", "updated_at", "status")

				dbMock.ExpectQuery("SELECT").
					WithArgs("email").
					WillReturnRows(rows)
			},
			input: paymentUser{
				UserEmail: "email",
			},
			expect: []payment{
				{1, 1, 10.5, "user_email", "currency", "created_at", "updated_at", "status"},
				{2, 2, 10.5, "user_email", "currency", "created_at", "updated_at", "status"},
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
			input: paymentUser{
				UserID: 1,
			},
			err: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.getPayments(
				context.TODO(),
				paymentUser{
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
				dbMock.ExpectExec("DELETE").
					WithArgs(1, statusSuccess, statusFailure).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:  1,
			expect: 1,
			err:    nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectExec("DELETE").
					WithArgs(1, statusSuccess, statusFailure).
					WillReturnError(errors.New("delete error"))
			},
			input: 1,
			err:   errors.New("delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.deletePayment(
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

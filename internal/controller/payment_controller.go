package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Neumann88/payment-api-emulator/internal/contract"
	"github.com/Neumann88/payment-api-emulator/internal/entity"
	"github.com/Neumann88/payment-api-emulator/pkg/loggin"
	"github.com/Neumann88/payment-api-emulator/pkg/util"
)

type PaymentController struct {
	usecase contract.PaymentUseCase
	logger  loggin.ILogger
}

func NewPaymentController(l loggin.ILogger, u contract.PaymentUseCase) *PaymentController {
	return &PaymentController{
		logger:  l,
		usecase: u,
	}
}

const (
	createPayment          = "/payment"
	updateStatusByID       = "/payments/{id}/status"
	getStatusByID          = "/payments/{id}/status"
	getPaymentsByUserEmail = "/payments/user" // query /payments/user?email=email
	getPaymentsByUserID    = "/payments/user/{id}"
	cancelPaymentByID      = "/payments/{id}"
)

func (c *PaymentController) Register(router *mux.Router) *mux.Router {
	router.HandleFunc(createPayment, c.createPayment).Methods(http.MethodPost)
	router.HandleFunc(updateStatusByID, c.updateStatus).Methods(http.MethodPut)
	router.HandleFunc(getStatusByID, c.getStatus).Methods(http.MethodGet)
	router.HandleFunc(getPaymentsByUserEmail, c.getPaymentsByUserEmail).Methods(http.MethodGet)
	router.HandleFunc(getPaymentsByUserID, c.getPaymentsByUserID).Methods(http.MethodGet)
	router.HandleFunc(cancelPaymentByID, c.cancelPayment).Methods(http.MethodPut)

	return router
}

func (c *PaymentController) createPayment(w http.ResponseWriter, r *http.Request) {
	var input entity.PaymentInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, util.InvalidBodyData, http.StatusBadRequest)

		return
	}

	if ok := util.IsEmail(input.UserEmail); !ok {
		http.Error(w, util.InvalidBodyEmail, http.StatusBadRequest)

		return
	}

	id, err := c.usecase.CreatePayment(
		r.Context(),
		input,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(
		entity.PaymentStatus{
			ID: id,
		},
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}
}

func (c *PaymentController) updateStatus(w http.ResponseWriter, r *http.Request) {
	paymentID, err := util.GetQueryID(r)

	if err != nil {
		http.Error(w, util.InvalidQueryID, http.StatusBadRequest)

		return
	}

	var input entity.PaymentStatus
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, util.InvalidBodyData, http.StatusBadRequest)

		return
	}

	input.ID = paymentID

	err = c.usecase.UpdateStatus(
		r.Context(),
		input,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *PaymentController) getStatus(w http.ResponseWriter, r *http.Request) {
	paymentID, err := util.GetQueryID(r)

	if err != nil {
		http.Error(w, util.InvalidQueryID, http.StatusBadRequest)

		return
	}

	status, err := c.usecase.GetStatus(
		r.Context(),
		paymentID,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(
		entity.PaymentStatus{
			Status: status,
		},
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}
}

func (c *PaymentController) getPaymentsByUserEmail(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("email")

	if ok := util.IsEmail(userEmail); !ok {
		http.Error(w, util.InvalidQueryEmail, http.StatusBadRequest)

		return
	}

	data, err := c.usecase.GetPayments(
		r.Context(),
		entity.PaymentUser{
			UserEmail: userEmail,
		},
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(
		entity.PaymentsData{
			Data: data,
		},
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}
}

func (c *PaymentController) getPaymentsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetQueryID(r)

	if err != nil {
		http.Error(w, util.InvalidQueryID, http.StatusBadRequest)

		return
	}

	data, err := c.usecase.GetPayments(
		r.Context(),
		entity.PaymentUser{
			UserID: userID,
		},
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(
		entity.PaymentsData{
			Data: data,
		},
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}
}

func (c *PaymentController) cancelPayment(w http.ResponseWriter, r *http.Request) {
	paymentID, err := util.GetQueryID(r)

	if err != nil {
		http.Error(w, util.InvalidQueryID, http.StatusBadRequest)

		return
	}

	err = c.usecase.CancelPayment(
		r.Context(),
		paymentID,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, util.InternalServerError, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

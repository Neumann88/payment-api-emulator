package payment

import (
	"encoding/json"
	"net/http"

	"github.com/Neumann88/payment-api-emulator/pkg/loggin"
	"github.com/gorilla/mux"
)

type controller struct {
	usecase paymentUseCase
	logger  loggin.ILogger
}

func NewPaymentController(l loggin.ILogger, u paymentUseCase) *controller {
	return &controller{
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
	deletePaymentByID      = "/payments/{id}"
)

func (c *controller) Register(router *mux.Router) *mux.Router {
	router.HandleFunc(createPayment, c.createPayment).Methods(http.MethodPost)
	router.HandleFunc(updateStatusByID, c.updateStatus).Methods(http.MethodPut)
	router.HandleFunc(getStatusByID, c.getStatus).Methods(http.MethodGet)
	router.HandleFunc(getPaymentsByUserEmail, c.getPaymentsByUserEmail).Methods(http.MethodGet)
	router.HandleFunc(getPaymentsByUserID, c.getPaymentsByUserID).Methods(http.MethodGet)
	router.HandleFunc(deletePaymentByID, c.deletePayment).Methods(http.MethodDelete)
	return router
}

func (c *controller) createPayment(w http.ResponseWriter, r *http.Request) {
	var input paymentInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, invalidBodyData, http.StatusBadRequest)
		return
	}

	if ok := isEmail(input.UserEmail); !ok {
		http.Error(w, invalidBodyEmail, http.StatusBadRequest)
		return
	}

	id, err := c.usecase.createPayment(
		r.Context(),
		input,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		paymentStatus{
			ID: id,
		},
	)
}

func (c *controller) updateStatus(w http.ResponseWriter, r *http.Request) {
	paymentID, err := getQueryId(r)

	if err != nil {
		http.Error(w, invalidQueryID, http.StatusBadRequest)
		return
	}

	var input paymentStatus
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, invalidBodyData, http.StatusBadRequest)
		return
	}

	input.ID = paymentID

	err = c.usecase.updateStatus(
		r.Context(),
		input,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *controller) getStatus(w http.ResponseWriter, r *http.Request) {
	paymentID, err := getQueryId(r)

	if err != nil {
		http.Error(w, invalidQueryID, http.StatusBadRequest)
		return
	}

	status, err := c.usecase.getStatus(
		r.Context(),
		paymentID,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		paymentStatus{
			Status: status,
		},
	)
}

func (c *controller) getPaymentsByUserEmail(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("email")
	if ok := isEmail(userEmail); !ok {
		http.Error(w, invalidQueryEmail, http.StatusBadRequest)
		return
	}

	data, err := c.usecase.getPayments(
		r.Context(),
		paymentUser{
			UserEmail: userEmail,
		},
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		paymentsData{
			Data: data,
		},
	)
}

func (c *controller) getPaymentsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := getQueryId(r)

	if err != nil {
		http.Error(w, invalidQueryID, http.StatusBadRequest)
		return
	}

	data, err := c.usecase.getPayments(
		r.Context(),
		paymentUser{
			UserID: userID,
		},
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		paymentsData{
			Data: data,
		},
	)
}

func (c *controller) deletePayment(w http.ResponseWriter, r *http.Request) {
	paymentID, err := getQueryId(r)

	if err != nil {
		http.Error(w, invalidQueryID, http.StatusBadRequest)
		return
	}

	err = c.usecase.deletePayment(
		r.Context(),
		paymentID,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

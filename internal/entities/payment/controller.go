package payment

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Controller struct {
	usecase PaymentUsecase
}

func NewPaymentController(u PaymentUsecase) *Controller {
	return &Controller{
		usecase: u,
	}
}

const (
	CREATE                     = "/payment"
	GET_STATUS_BY_ID           = "/payment/{id}/status"
	UPDATE_STATUS_BY_ID        = "/payment/{id}/status"
	GET_PAYMENTS_BY_USER_ID    = "/payments" // query /payments?user_id="..."
	GET_PAYMENTS_BY_USER_EMAIL = "/payments" // query /payments?user_email="..."
	CANCEL_BY_ID               = "/payment/{id}"
)

func (c *Controller) Register(router *mux.Router) {
	router.HandleFunc(CREATE, c.createPayment).Methods(http.MethodPost)
	// router.HandleFunc(GET_STATUS_BY_ID, s.getStatus).Methods(http.MethodGet)
	// router.HandleFunc(UPDATE_STATUS_BY_ID, s.updateStatus).Methods(http.MethodPut)
	// router.HandleFunc(GET_PAYMENTS_BY_USER_ID, s.getPaymentsByUserID).Methods(http.MethodGet)
	// router.HandleFunc(GET_PAYMENTS_BY_USER_EMAIL, s.getPaymentsByUserEmail).Methods(http.MethodGet)
	// router.HandleFunc(CANCEL_BY_ID, s.cancelPayment).Methods(http.MethodDelete)
}

func (c *Controller) createPayment(w http.ResponseWriter, r *http.Request) {
	var input Payment
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := c.usecase.createPayment(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// func (p *PaymentHandler) updateStatus(w http.ResponseWriter, r *http.Request) {
// 	var input AccoutAddBalanceRequest
// 	err := json.NewDecoder(r.Body).Decode(&input)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	err = a.service.UpdateBalance(input.UserID, input.Amount)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// }

// func (p *PaymentHandler) getStatus(w http.ResponseWriter, r *http.Request) {
// 	var input AccoutAddBalanceRequest
// 	err := json.NewDecoder(r.Body).Decode(&input)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	err = a.service.UpdateBalance(input.UserID, input.Amount)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// }

// func (p *PaymentHandler) getPaymentsByUserEmail(w http.ResponseWriter, r *http.Request) {
// 	var input AccoutAddBalanceRequest
// 	err := json.NewDecoder(r.Body).Decode(&input)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	err = a.service.UpdateBalance(input.UserID, input.Amount)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// }

// func (p *PaymentHandler) getPaymentsByUserID(w http.ResponseWriter, r *http.Request) {
// 	var input AccoutAddBalanceRequest
// 	err := json.NewDecoder(r.Body).Decode(&input)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	err = a.service.UpdateBalance(input.UserID, input.Amount)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// }

// func (p *PaymentHandler) cancelPayment(w http.ResponseWriter, r *http.Request) {
// 	var input AccoutAddBalanceRequest
// 	err := json.NewDecoder(r.Body).Decode(&input)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	err = a.service.UpdateBalance(input.UserID, input.Amount)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// }

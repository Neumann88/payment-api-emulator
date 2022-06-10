package controller

import (
	"github.com/Neumann88/payment-api-emulator/internal/entities/payment"
	"github.com/Neumann88/payment-api-emulator/internal/usecase"
	"github.com/gorilla/mux"
)

type Controller struct {
	PaymentController payment.PaymentController
}

func NewController(s *usecase.Usecase) *Controller {
	return &Controller{
		PaymentController: payment.NewPaymentController(s.PaymentUsecase),
	}
}

func (h *Controller) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	h.PaymentController.Register(router)

	return router
}

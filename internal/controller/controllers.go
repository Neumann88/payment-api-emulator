package controller

import (
	"github.com/Neumann88/payment-api-emulator/internal/entities/payment"
	"github.com/Neumann88/payment-api-emulator/internal/usecase"
	"github.com/Neumann88/payment-api-emulator/pkg/loggin"
	"github.com/gorilla/mux"
)

type Controller struct {
	PaymentController payment.PaymentController
}

func NewController(l loggin.ILogger, s *usecase.Usecase) *Controller {
	return &Controller{
		PaymentController: payment.NewPaymentController(l, s.PaymentUsecase),
	}
}

func (h *Controller) InitRoutes() *mux.Router {
	router := mux.NewRouter()
	h.PaymentController.Register(router)
	return router
}

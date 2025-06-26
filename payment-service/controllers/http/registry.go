package controllers

import (
	controllers "github.com/anddriii/kita-futsal/payment-service/controllers/http/payment"
	"github.com/anddriii/kita-futsal/payment-service/service"
)

type Registry struct {
	service service.IServiceRegistry
}

// GetPayment implements IControllerRegistry.
func (r *Registry) GetPayment() controllers.IPaymentController {
	return controllers.NewPaymentController(r.service)
}

type IControllerRegistry interface {
	GetPayment() controllers.IPaymentController
}

func NewControllerRegistry(service service.IServiceRegistry) IControllerRegistry {
	return &Registry{
		service: service,
	}
}

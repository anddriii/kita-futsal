package controllers

import (
	controllers "github.com/anddriii/kita-futsal/order-service/controllers/http/order"
	"github.com/anddriii/kita-futsal/order-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetOrder() controllers.IOrderController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

func (r *Registry) GetOrder() controllers.IOrderController {
	return controllers.NewOrderController(r.service)
}

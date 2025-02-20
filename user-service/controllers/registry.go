package controllers

import (
	controllers "github.com/anddriii/kita-futsal/user-service/controllers/user"
	"github.com/anddriii/kita-futsal/user-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetUserController() controllers.IUserController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

// GetUserController implements IControllerRegistry.
func (r *Registry) GetUserController() controllers.IUserController {
	return controllers.NewUserController(r.service)
}

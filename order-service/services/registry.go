package services

import (
	"github.com/anddriii/kita-futsal/order-service/clients"
	"github.com/anddriii/kita-futsal/order-service/repositories"
	services "github.com/anddriii/kita-futsal/order-service/services/order"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
	client     clients.IClientRegistry
}

type IServiceRegistry interface {
	GetOrder() services.IOrdderService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry, client clients.IClientRegistry) IServiceRegistry {
	return &Registry{repository: repository, client: client}
}

func (r *Registry) GetOrder() services.IOrdderService {
	return services.NewOrderService(r.repository, r.client)
}

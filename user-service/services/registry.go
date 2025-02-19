package services

import (
	"github.com/anddriii/kita-futsal/user-service/repositories"
	service "github.com/anddriii/kita-futsal/user-service/services/user"
)

type Registry struct {
	repository repositories.IRepoRegistry
}

type IServiceRegistry interface {
	GetUser() service.IUserService
}

func NewServiceRegistry(repository repositories.IRepoRegistry) IServiceRegistry {
	return &Registry{repository: repository}
}

// GetUser implements IServiceRegistry.
func (r *Registry) GetUser() service.IUserService {
	return service.NewUserService(r.repository)
}

package clients

import (
	"github.com/anddriii/kita-futsal/field-service/clients/config"
	clients "github.com/anddriii/kita-futsal/field-service/clients/user"
	config2 "github.com/anddriii/kita-futsal/field-service/config"
)

type ClientRegistry struct {
}

type IClientRegistry interface {
	GetUser() clients.IUserClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}

func (c *ClientRegistry) GetUser() clients.IUserClient {
	return clients.NewUserClient(
		config.NewClientConfig(
			config.WithBaseURL(config2.Config.InternalService.User.Host),
			config.WithSignatureKey(config2.Config.InternalService.User.SignatureKey),
		))

}

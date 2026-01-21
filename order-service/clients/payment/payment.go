package clients

import (
	"context"

	"github.com/anddriii/kita-futsal/order-service/clients/config"
	"github.com/anddriii/kita-futsal/order-service/domain/dto"
	"github.com/google/uuid"
)

type PaymentClient struct {
	client config.IClientConfig
}

// CreatePaymentLink implements IPaymentClient.
func (p *PaymentClient) CreatePaymentLink(context.Context, *dto.PaymentRequest) (*PaymentData, error) {
	panic("unimplemented")
}

// GetPaymentByUUID implements IPaymentClient.
func (p *PaymentClient) GetPaymentByUUID(context.Context, uuid.UUID) (*PaymentData, error) {
	panic("unimplemented")
}

type IPaymentClient interface {
	GetPaymentByUUID(context.Context, uuid.UUID) (*PaymentData, error)
	CreatePaymentLink(context.Context, *dto.PaymentRequest) (*PaymentData, error)
}

func NewPaymentClient(client config.IClientConfig) IPaymentClient {
	return &PaymentClient{client: client}
}

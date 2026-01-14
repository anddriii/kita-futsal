package clients

import (
	"context"

	"github.com/anddriii/kita-futsal/order-service/clients/config"
	"github.com/google/uuid"
)

type FieldClient struct {
	client config.IClientConfig
}

type IFieldClient interface {
	GetFieldByUUID(context.Context, uuid.UUID) (*FieldData, error)
	UpdateStatus()
}

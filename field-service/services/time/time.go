package time

import (
	"context"

	"github.com/anddriii/kita-futsal/field-service/domains/dto"
)

type ITimeService interface {
	GetAll(ctx context.Context) ([]dto.TimeResponse, error)
	GetByUUID(ctx context.Context, uuid string) (*dto.TimeResponse, error)
	Create(ctx context.Context, req *dto.TimeRequest) (*dto.TimeResponse, error)
}

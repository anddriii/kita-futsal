package services

import (
	"context"

	"github.com/anddriii/kita-futsal/order-service/common/util"
	"github.com/anddriii/kita-futsal/order-service/domain/dto"
)

type IOrdderService interface {
	GetAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) (*util.PaginationResult, error)
	GetByUUID(ctx context.Context, uuid string) (*dto.OrderResponse, error)
	GetOrderByUserID(ctx context.Context) ([]dto.OrderByUserIDResponse, error)
	Create(ctx context.Context, req *dto.OrderRequest) (*dto.OrderResponse, error)
	HandlePayment(ctx context.Context, request *dto.PaymentData) error
}

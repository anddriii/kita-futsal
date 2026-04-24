package service

import (
	"context"

	"github.com/anddriii/kita-futsal/payment-service/common/util"
	"github.com/anddriii/kita-futsal/payment-service/domains/dto"
)

type IPaymentService interface {
	GetAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) (*util.PaginationResult, error)
	GetByUUID(ctx context.Context, uuid string) (*dto.PaymentResponse, error)
	Create(ctx context.Context, req *dto.PaymentRequest) (*dto.PaymentResponse, error)
	WebHook(ctx context.Context, req *dto.Webhook) error
}

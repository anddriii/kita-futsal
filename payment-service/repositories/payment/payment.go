package repositories

import (
	"context"

	"github.com/anddriii/kita-futsal/payment-service/domains/dto"
	"github.com/anddriii/kita-futsal/payment-service/domains/models"
	"gorm.io/gorm"
)

type IPaymentRepository interface {
	FindAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) ([]models.Payment, int64, error)
	FindByUUID(ctx context.Context, uuid string) (*models.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*models.Payment, error)
	Create(ctx context.Context, db *gorm.DB, req *dto.PaymentRequest) (*models.Payment, error)
	Update(ctx context.Context, db *gorm.DB, orderID string, req *dto.UpdatePaymentRequest) (*models.Payment, error)
}

package order

import (
	"context"

	"github.com/anddriii/kita-futsal/order-service/domain/dto"
	"github.com/anddriii/kita-futsal/order-service/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	FindAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) ([]models.Order, int64, error)
	FindByUserID(ctx context.Context, userID string) ([]models.Order, error)
	FindByUUID(ctx context.Context, uuid string) (*models.Order, error)
	Create(ctx context.Context, tx *gorm.DB, param *models.Order) (*models.Order, error)
	Update(ctx context.Context, tx *gorm.DB, request *models.Order, uuid uuid.UUID) error
}

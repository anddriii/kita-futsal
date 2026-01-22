package order

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/anddriii/kita-futsal/order-service/domain/dto"
	"github.com/anddriii/kita-futsal/order-service/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func (o *OrderRepository) incrementCode(ctx context.Context) (*string, error) {
	var (
		order  *models.Order
		result string
		today  = time.Now().Format("20060102")
	)

	err := o.db.WithContext(ctx).Order("id desc").First(&order).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	if order.ID != 0 {
		orderCode := order.Code
		splitOrderName, _ := strconv.Atoi(orderCode[4:9])
		code := splitOrderName + 1
		result = fmt.Sprintf("ORD-%5d-%s", code, today)
	} else {
		result = fmt.Sprintf("ORD-%5d-%s", 1, today)
	}

	return &result, nil
}

// Create implements IOrderRepository.
func (o OrderRepository) Create(ctx context.Context, tx *gorm.DB, param *models.Order) (*models.Order, error) {
	panic("unimplemented")
}

// FindAllWithPagination implements IOrderRepository.
func (o OrderRepository) FindAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) ([]models.Order, int64, error) {
	panic("unimplemented")
}

// FindByUUID implements IOrderRepository.
func (o OrderRepository) FindByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	panic("unimplemented")
}

// FindByUserID implements IOrderRepository.
func (o OrderRepository) FindByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	panic("unimplemented")
}

// Update implements IOrderRepository.
func (o OrderRepository) Update(ctx context.Context, tx *gorm.DB, request *models.Order, uuid uuid.UUID) error {
	panic("unimplemented")
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return OrderRepository{db: db}
}

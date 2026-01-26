package order

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	errWrap "github.com/anddriii/kita-futsal/order-service/common/error"
	errConstant "github.com/anddriii/kita-futsal/order-service/constants/error"
	errOrder "github.com/anddriii/kita-futsal/order-service/constants/error/order"
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
	code, err := o.incrementCode(ctx)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		UUID:   uuid.New(),
		Code:   *code,
		UserID: param.UserID,
		Amount: param.Amount,
		Date:   param.Date,
		Status: param.Status,
		IsPaid: param.IsPaid,
	}

	err = tx.
		WithContext(ctx).
		Create(order).
		Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return order, nil
}

// FindAllWithPagination implements IOrderRepository.
func (o OrderRepository) FindAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) ([]models.Order, int64, error) {
	var (
		orders []models.Order
		sort   string
		total  int64
	)
	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := o.db.
		WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&orders).
		Error
	if err != nil {
		return nil, 0, err
	}

	err = o.db.
		WithContext(ctx).
		Model(&models.Order{}).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return orders, total, nil
}

// FindByUUID implements IOrderRepository.
func (o OrderRepository) FindByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	var order models.Order
	err := o.db.
		WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&order).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errOrder.ErrOrderNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &order, nil
}

// FindByUserID implements IOrderRepository.
func (o OrderRepository) FindByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order
	err := o.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&orders).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errOrder.ErrOrderNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return orders, nil
}

// Update implements IOrderRepository.
func (o OrderRepository) Update(ctx context.Context, tx *gorm.DB, request *models.Order, uuid uuid.UUID) error {
	err := tx.
		WithContext(ctx).
		Model(&models.Order{}).
		Where("uuid = ?", uuid).
		Updates(request).
		Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return OrderRepository{db: db}
}

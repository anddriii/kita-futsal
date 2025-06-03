package repositories

import (
	"context"

	errWrap "github.com/anddriii/kita-futsal/payment-service/common/error"
	errConst "github.com/anddriii/kita-futsal/payment-service/constants/error"
	"github.com/anddriii/kita-futsal/payment-service/domains/dto"
	"github.com/anddriii/kita-futsal/payment-service/domains/models"
	"gorm.io/gorm"
)

type PaymentHistoryRepository struct {
	db *gorm.DB
}

type IPaymentHistoryRepository interface {
	Create(ctx context.Context, tx *gorm.DB, req *dto.PaymentHistoryRequest) error
}

// Create implements IPaymentHistoryRepository.
func (p *PaymentHistoryRepository) Create(ctx context.Context, tx *gorm.DB, req *dto.PaymentHistoryRequest) error {
	paymentHistory := models.PaymentHistory{
		PaymentID: req.PaymentID,
		Status:    req.Status,
	}

	err := tx.WithContext(ctx).Create(&paymentHistory).Error
	if err != nil {
		return errWrap.WrapError(errConst.ErrSQLError)
	}

	return nil
}

func NewPaymentHistoryRepository(db *gorm.DB) IPaymentHistoryRepository {
	return &PaymentHistoryRepository{db: db}
}

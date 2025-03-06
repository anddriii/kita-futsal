package repositories

import (
	"context"
	"errors"
	"fmt"

	errWrap "github.com/anddriii/kita-futsal/field-service/common/error"
	errConst "github.com/anddriii/kita-futsal/field-service/constants/error"
	"github.com/anddriii/kita-futsal/field-service/domains/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeRepository struct {
	db *gorm.DB
}

// Create implements ITimeRepository.
func (t *TimeRepository) Create(ctx context.Context, time *models.Time) (*models.Time, error) {
	time.UUID = uuid.New()
	fmt.Println("time", time)
	err := t.db.WithContext(ctx).Create(time).Error
	if err != nil {
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return time, nil
}

// FindAll implements ITimeRepository.
func (t *TimeRepository) FindAll(ctx context.Context) ([]models.Time, error) {
	var times []models.Time
	err := t.db.WithContext(ctx).Find(&times).Error
	if err != nil {
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}
	return times, nil
}

// FindById implements ITimeRepository.
func (t *TimeRepository) FindById(ctx context.Context, id int) (*models.Time, error) {
	var time models.Time
	err := t.db.WithContext(ctx).Where("id = ?", id).First(&time).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return &time, nil
}

// FindByUUID implements ITimeRepository.
func (t *TimeRepository) FindByUUID(ctx context.Context, uuid string) (*models.Time, error) {
	var time models.Time
	err := t.db.WithContext(ctx).Where("uuid = ?", uuid).First(&time).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return &time, nil
}

func NewTimeRepository(db *gorm.DB) ITimeRepository {
	return &TimeRepository{db: db}
}

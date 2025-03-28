package repositoroes

import (
	"context"
	"errors"
	"fmt"

	errWrap "github.com/anddriii/kita-futsal/field-service/common/error"
	errConst "github.com/anddriii/kita-futsal/field-service/constants/error"
	errField "github.com/anddriii/kita-futsal/field-service/constants/error/field"
	"github.com/anddriii/kita-futsal/field-service/domains/dto"
	"github.com/anddriii/kita-futsal/field-service/domains/models"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldRepository struct {
	db *gorm.DB
}

func NewFieldRepository(db *gorm.DB) IFieldRepository {
	return &FieldRepository{db: db}
}

// Create implements IFieldRepository.
func (f *FieldRepository) Create(ctx context.Context, req *models.Field) (*models.Field, error) {
	field := models.Field{
		UUID:         uuid.New(),
		Code:         req.Code,
		Name:         req.Name,
		Image:        req.Image,
		PricePerHour: req.PricePerHour,
	}

	fmt.Print("sudah masuk ke database")

	err := f.db.WithContext(ctx).Create(field).Error
	if err != nil {
		log.Errorf("error from repositories", err)
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	fmt.Print("sudah dibuat ke database")

	return &field, nil
}

// Delete implements IFieldRepository.
func (f *FieldRepository) Delete(ctx context.Context, uuid string) error {
	err := f.db.WithContext(ctx).Where(uuid).Delete(&models.Field{}).Error
	if err != nil {
		return errWrap.WrapError(errConst.ErrSQLError)
	}

	return nil
}

// FindALlWithPagination implements IFieldRepository.
func (f *FieldRepository) FindALlWithPagination(ctx context.Context, param *dto.FieldRequestParam) ([]models.Field, int64, error) {
	var (
		fields []models.Field
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
	err := f.db.WithContext(ctx).Limit(limit).Offset(offset).Order(sort).Find(&fields).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConst.ErrSQLError)
	}

	err = f.db.WithContext(ctx).Model(&fields).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConst.ErrSQLError)
	}

	return fields, total, nil
}

// FindAllWithoutPagination implements IFieldRepository.
func (f *FieldRepository) FindAllWithoutPagination(ctx context.Context) ([]models.Field, error) {
	var fileds []models.Field
	err := f.db.WithContext(ctx).Find(fileds).Error
	if err != nil {
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return fileds, nil
}

// FindByUUID implements IFieldRepository.
func (f *FieldRepository) FindByUUID(ctx context.Context, uuid string) (*models.Field, error) {
	var fields models.Field
	err := f.db.WithContext(ctx).Where(uuid).First(&fields).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errField.ErrFieldNotFound)
		}
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return &fields, nil
}

// Update implements IFieldRepository.
func (f *FieldRepository) Update(ctx context.Context, uuid string, req *models.Field) (*models.Field, error) {
	field := models.Field{
		Code:         req.Code,
		Name:         req.Name,
		Image:        req.Image,
		PricePerHour: req.PricePerHour,
	}

	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&field).Error
	if err != nil {
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return &field, nil
}

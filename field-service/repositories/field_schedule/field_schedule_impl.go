package repositories

import (
	"context"
	"errors"
	"fmt"

	errWrap "github.com/anddriii/kita-futsal/field-service/common/error"
	"github.com/anddriii/kita-futsal/field-service/constants"
	errConst "github.com/anddriii/kita-futsal/field-service/constants/error"
	errField "github.com/anddriii/kita-futsal/field-service/constants/error/field_schedule"
	"github.com/anddriii/kita-futsal/field-service/domains/dto"
	"github.com/anddriii/kita-futsal/field-service/domains/models"
	"gorm.io/gorm"
)

type FieldScheduleRepository struct {
	db *gorm.DB
}

func NewFieldScheduleRepository(db *gorm.DB) IFieldScheduleRepository {
	return &FieldScheduleRepository{db: db}
}

// Create implements IFieldScheduleRepository.
func (f *FieldScheduleRepository) Create(ctx context.Context, req []models.FieldSchedule) error {
	err := f.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return errWrap.WrapError(errConst.ErrSQLError)
	}

	return nil
}

// Delete implements IFieldScheduleRepository.
func (f *FieldScheduleRepository) Delete(ctx context.Context, uuid string) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.FieldSchedule{}).Error
	if err != nil {
		return errWrap.WrapError(errConst.ErrSQLError)
	}

	return nil
}

// FindAllByIdAndDate implements IFieldScheduleRepository.
func (f *FieldScheduleRepository) FindAllByIdAndDate(ctx context.Context, FieldId int, date string) ([]models.FieldSchedule, error) {
	var fieldSchedules []models.FieldSchedule

	err := f.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("field_id = ?", FieldId).
		Where("date = ?", date).
		Joins("LEFT JOIN times ON field_schedules.time_id = times.id"). //menghubungkan tabel field_schedules dengan tabel times berdasarkan time_id.
		Order("times.start_time asc").                                  // Mengurutkan hasil berdasarkan start_time secara ascending.
		Find(&fieldSchedules).Error
	if err != nil {
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return fieldSchedules, nil
}

// FindAllWithPagination implements IFieldScheduleRepository.
func (f *FieldScheduleRepository) FindAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error) {
	var (
		fieldSchedules []models.FieldSchedule
		sort           string
		total          int64
	)

	// Jika user memberikan kolom sorting (SortColumn) dan arah sorting (SortOrder), maka formatnya dibuat sesuai input.
	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	/*
		limit → Menentukan jumlah data per halaman.
		offset → Menentukan dari data keberapa query akan dimulai.
				(Misalnya, Page = 2 dengan Limit = 10, maka offset = (2-1) * 10 = 10, artinya query akan mulai dari data ke-11.)
	*/
	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := f.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Limit(limit).   // Mengatur jumlah data yang diambil dalam satu halaman.
		Offset(offset). // Menentukan titik mulai pengambilan data berdasarkan halaman
		Order(sort).    // Mengurutkan hasil berdasarkan parameter sorting.
		Find(&fieldSchedules).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConst.ErrSQLError)
	}

	// menghitung jumlah total data yang tersedia tanpa pagination.
	err = f.db.WithContext(ctx).Model(&models.FieldSchedule{}).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConst.ErrSQLError)
	}
	return fieldSchedules, total, nil
}

// FindByDateAndTimeId implements IFieldScheduleRepository.
func (f *FieldScheduleRepository) FindByDateAndTimeId(ctx context.Context, date string, timeID int, fieldID int) (*models.FieldSchedule, error) {
	var fieldSchedule models.FieldSchedule
	err := f.db.WithContext(ctx).Where("date = ?", date).Where("time_id = ?", timeID).Where("field_id = ?", fieldID).First(&fieldSchedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}
	fmt.Println("response dari repo", fieldSchedule)

	return &fieldSchedule, nil
}

// FindByUUID implements IFieldScheduleRepository.
func (f *FieldScheduleRepository) FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error) {
	var fieldSchedule models.FieldSchedule

	err := f.db.WithContext(ctx).Preload("Field").Preload("Time").Where("uuid = ?", uuid).First(&fieldSchedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errField.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}
	return &fieldSchedule, nil
}

// Update implements IFieldScheduleRepository.
func (f *FieldScheduleRepository) Update(ctx context.Context, uuid string, req *models.FieldSchedule) (*models.FieldSchedule, error) {
	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedule.Date = req.Date
	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return fieldSchedule, nil
}

// UpdateStatus implements IFieldScheduleRepository.
func (f *FieldScheduleRepository) UpdateStatus(ctx context.Context, status constants.FieldScheduleStatus, uuid string) error {
	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	fieldSchedule.Status = status
	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		return errWrap.WrapError(errConst.ErrSQLError)
	}

	return nil
}

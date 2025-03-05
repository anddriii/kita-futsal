package repositories

import (
	"context"

	"github.com/anddriii/kita-futsal/field-service/constants"
	"github.com/anddriii/kita-futsal/field-service/domains/dto"
	"github.com/anddriii/kita-futsal/field-service/domains/models"
)

type IFieldScheduleRepository interface {
	FindAllWithPagination(ctx context.Context, req *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error)
	FindAllByIdAndDate(ctx context.Context, FieldId int, date string) ([]models.FieldSchedule, error)
	FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error)
	FindByDateAndTimeId(ctx context.Context, date string, timeID int, fieldID int) (*models.FieldSchedule, error)
	Create(ctx context.Context, req []models.FieldSchedule) error
	Update(ctx context.Context, uuid string, req *models.FieldSchedule) (*models.FieldSchedule, error)
	UpdateStatus(ctx context.Context, status constants.FieldScheduleStatus, uuid string) error
	Delete(ctx context.Context, uuid string) error
}

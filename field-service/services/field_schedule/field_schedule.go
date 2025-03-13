package services

import (
	"context"

	"github.com/anddriii/kita-futsal/field-service/common/util"
	"github.com/anddriii/kita-futsal/field-service/domains/dto"
)

type IFieldScheduleService interface {
	FindAllWithPagination(ctc context.Context, req *dto.FieldScheduleRequestParam) (*util.PaginationResult, error)
	FindAllFieldByIdAndDate(ctx context.Context, uuid string, date string) ([]dto.FieldScheduleForBookingReponse, error)
	FindByUUID(ctx context.Context, uuid string) (dto.FieldScheduleResponse, error)
	GenereateScheduleForOneMonth(ctx context.Context, req dto.GenerateFieldScheduleForOneMonthRequest) error
	Create(ctx context.Context, req *dto.FieldScheduleRequest) error
	Update(ctx context.Context, uuid string, req *dto.UpdateFieldRequest) (*dto.FieldScheduleResponse, error)
	UpdateStatus(ctx context.Context, req dto.UpdateStatusFieldScheduleRequest)
	Delete(ctx context.Context, uuid string) error
}

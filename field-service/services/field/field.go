package services

import (
	"context"

	"github.com/anddriii/kita-futsal/field-service/common/util"
	"github.com/anddriii/kita-futsal/field-service/domains/dto"
)

type IFieldService interface {
	GetAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) (*util.PaginationResult, error)
	GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error)
	GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error)
	Create(ctx context.Context, req *dto.FieldRequest) (*dto.FieldResponse, error)
	Update(ctx context.Context, uuid string, req *dto.UpdateFieldRequest) (*dto.FieldResponse, error)
	Delete(ctx context.Context, uuid string) error
}

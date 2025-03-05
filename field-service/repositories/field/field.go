package repositoroes

import (
	"context"

	"github.com/anddriii/kita-futsal/field-service/domains/dto"
	"github.com/anddriii/kita-futsal/field-service/domains/models"
)

type IFieldRepository interface {
	FindALlWithPagination(ctx context.Context, req *dto.FieldRequestParam) ([]models.Field, int64, error)
	FindAllWithoutPagination(ctx context.Context) ([]models.Field, error)
	FindByUUID(ctx context.Context, uuid string) (*models.Field, error)
	Create(ctx context.Context, req *models.Field) (*models.Field, error)
	Update(ctx context.Context, uuid string, req *models.Field) (*models.Field, error)
	Delete(ctx context.Context, uuid string) error
}

package repositories

import (
	"context"

	"github.com/anddriii/kita-futsal/field-service/domains/models"
)

type ITimeRepository interface {
	FindAll(ctx context.Context) ([]models.Time, error)
	FindByUUID(ctx context.Context, uuid string) (*models.Time, error)
	FindById(ctx context.Context, id int) (*models.Time, error)
	Create(ctx context.Context, time *models.Time) (*models.Time, error)
}

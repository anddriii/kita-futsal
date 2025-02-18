package repositories

import (
	"context"

	"github.com/anddriii/kita-futsal/user-service/domain/dto"
	"github.com/anddriii/kita-futsal/user-service/domain/models"
)

type UserRepo interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error)
	Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
}

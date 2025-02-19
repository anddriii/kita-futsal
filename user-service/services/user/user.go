package services

import (
	"context"

	"github.com/anddriii/kita-futsal/user-service/domain/dto"
)

type IUserService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(ctx context.Context, req *dto.UpdateRequest, username string) (*dto.UserResponse, error)
	GetUserLogin(ctx context.Context) (*dto.UserResponse, error)
	GetUserUUID(ctx context.Context, uuid string) (*dto.UserResponse, error)
	ifUsernameExist(ctx context.Context, username string) bool
	ifEmailExist(ctx context.Context, email string) bool
}

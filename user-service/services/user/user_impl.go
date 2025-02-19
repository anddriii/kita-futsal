package services

import (
	"context"
	"strings"
	"time"

	"github.com/anddriii/kita-futsal/user-service/config"
	"github.com/anddriii/kita-futsal/user-service/domain/dto"
	"github.com/anddriii/kita-futsal/user-service/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRepoRegistry
}

func NewUserService(repository repositories.IRepoRegistry) IUserService {
	return &UserService{repository: repository}
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

// Login implements IUserService.
func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginReponse, error) {
	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpireTime) * time.Minute).Unix()

	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        strings.ToLower(user.Role.Code),
	}

	claims := Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	tokenString, err := token.SignedString([]byte(config.Config.JwtScretKey))
	if err != nil {
		return nil, err
	}

	response := dto.LoginReponse{
		User:  *data,
		Token: tokenString,
	}

	return &response, nil
}

// Register implements IUserService.
func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	panic("unimplemented")
}

// Update implements IUserService.
func (u *UserService) Update(ctx context.Context, req *dto.UpdateRequest, username string) (*dto.UserResponse, error) {
	panic("unimplemented")
}

// GetUserLogin implements IUserService.
func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	panic("unimplemented")
}

// GetUserUUID implements IUserService.
func (u *UserService) GetUserUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	panic("unimplemented")
}

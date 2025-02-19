package services

import (
	"context"
	"strings"
	"time"

	"github.com/anddriii/kita-futsal/user-service/config"
	"github.com/anddriii/kita-futsal/user-service/constants"
	errConst "github.com/anddriii/kita-futsal/user-service/constants/error"
	"github.com/anddriii/kita-futsal/user-service/domain/dto"
	"github.com/anddriii/kita-futsal/user-service/domain/models"
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
func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	//Verifikasi Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	//Menentukan Waktu Kadaluarsa Token
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	response := dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return &response, nil
}

func (u *UserService) ifUsernameExist(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByUsername(ctx, username)
	if err != nil {
		return false
	}

	//cek apakah username pernah digunakan atau belum
	if user != nil {
		return true
	}

	return false
}

// ifEmailExist implements IUserService.
func (u *UserService) ifEmailExist(ctx context.Context, email string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	//cek apakah email pernah digunakan atau belum
	if user != nil {
		return true
	}

	return false
}

// Register implements IUserService.
func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if u.ifUsernameExist(ctx, req.Username) {
		return nil, errConst.ErrUsernameExist
	}

	if u.ifEmailExist(ctx, req.Email) {
		return nil, errConst.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, errConst.ErrPasswordDoesNotMatch
	}

	reqUser := &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    string(hashedPW),
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		RoleId:      constants.User,
	}

	user, err := u.repository.GetUser().Register(ctx, reqUser)
	if err != nil {
		return nil, err
	}

	response := dto.RegisterResponse{
		User: *&dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
	}

	return &response, nil
}

// Update implements IUserService.
func (u *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	var (
		password                  string
		checkUsername, checkEmail *models.User
		hashedPW                  []byte
		user, userResult          *models.User
		err                       error
		data                      dto.UserResponse
	)

	user, err = u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	IsUsernameExist := u.ifUsernameExist(ctx, req.Username)
	if IsUsernameExist && user.Username != req.Username {
		checkUsername, err = u.repository.GetUser().FindByUsername(ctx, req.Username)
		if err != nil {
			return nil, err
		}

		if checkUsername != nil {
			return nil, errConst.ErrUsernameExist
		}
	}

	IsEmailExist := u.ifEmailExist(ctx, req.Email)
	if IsEmailExist && user.Email != req.Email {
		checkEmail, err = u.repository.GetUser().FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}

		if checkEmail != nil {
			return nil, errConst.ErrEmailExist
		}
	}

	if req.Password != nil {
		if *req.Password != *req.ConfirmPassword {
			return nil, errConst.ErrPasswordDoesNotMatch
		}
		hashedPW, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		password = string(hashedPW)
	}

	userReq := &dto.UpdateRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    &password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	userResult, err = u.repository.GetUser().Update(ctx, userReq, uuid)
	if err != nil {
		return nil, err
	}

	data = dto.UserResponse{
		UUID:        userResult.UUID,
		Name:        userResult.Name,
		Username:    userResult.Username,
		Email:       userResult.Email,
		PhoneNumber: userResult.PhoneNumber,
	}

	return &data, nil

}

// GetUserLogin implements IUserService.
func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		//mengambil informasi pengguna yang sedang login dari context
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		Email:       userLogin.Email,
		Role:        userLogin.Role,
		PhoneNumber: userLogin.PhoneNumber,
	}

	return &data, nil
}

// GetUserUUID implements IUserService.
func (u *UserService) GetUserUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	data := dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}

	return &data, nil
}

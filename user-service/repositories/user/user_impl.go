package repositories

import (
	"context"
	"errors"
	"log"

	errWrap "github.com/anddriii/kita-futsal/user-service/common/error"
	errConstant "github.com/anddriii/kita-futsal/user-service/constants/error"
	"github.com/anddriii/kita-futsal/user-service/domain/dto"
	"github.com/anddriii/kita-futsal/user-service/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	return &UserRepoImpl{db: db}
}

// Register implements UserRepo.
func (u *UserRepoImpl) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:        uuid.New(),
		Name:        req.Name,
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleId:      req.RoleId,
	}

	err := u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		log.Println("Error inserting user to DB:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError) // jika terjadi error akan memanipulasi log error agar tidak menampilkan "Query error" nya
	}

	return &user, nil
}

// Update implements UserRepo.
func (u *UserRepoImpl) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*models.User, error) {
	user := models.User{
		Name:        req.Name,
		Username:    req.Username,
		Password:    *req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	err := u.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&user).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil

}

// FindByEmail implements UserRepo.
func (u *UserRepoImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := u.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

// FindByUUID implements UserRepo.
func (u *UserRepoImpl) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User

	err := u.db.WithContext(ctx).Preload("Role").Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

// FindByUsername implements UserRepo.
func (u *UserRepoImpl) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := u.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		log.Println("error repositories:", err)
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

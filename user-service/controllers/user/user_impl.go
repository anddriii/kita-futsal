package controllers

import (
	"net/http"

	errWrap "github.com/anddriii/kita-futsal/user-service/common/error"
	"github.com/anddriii/kita-futsal/user-service/common/response"
	"github.com/anddriii/kita-futsal/user-service/domain/dto"
	"github.com/anddriii/kita-futsal/user-service/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserControllers struct {
	UserService services.IServiceRegistry
}

func NewUserController(userService services.IServiceRegistry) IUserController {
	return &UserControllers{UserService: userService}
}

// Login implements IUserController.
func (u *UserControllers) Login(ctx *gin.Context) {
	request := dto.LoginRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.WrapError(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
	}

	user, err := u.UserService.GetUser().Login(ctx, &request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code:  http.StatusOK,
		Data:  user.User,
		Token: &user.Token,
		Gin:   ctx,
	})
}

// Register implements IUserController.
func (u *UserControllers) Register(ctx *gin.Context) {
	request := dto.RegisterRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.WrapError(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
	}

	user, err := u.UserService.GetUser().Register(ctx, &request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user.User,
		Gin:  ctx,
	})
}

// Update implements IUserController.
func (u *UserControllers) Update(ctx *gin.Context) {
	request := dto.UpdateRequest{}
	uuid := ctx.Param("uuid")

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.WrapError(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
	}

	user, err := u.UserService.GetUser().Update(ctx, &request, uuid)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

// GetUserUUID implements IUserController.
func (u *UserControllers) GetUserUUID(ctx *gin.Context) {
	user, err := u.UserService.GetUser().GetUserLogin(ctx.Request.Context())
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

// GetUserLogin implements IUserController.
func (u *UserControllers) GetUserLogin(ctx *gin.Context) {
	user, err := u.UserService.GetUser().GetUserUUID(ctx.Request.Context(), ctx.Param("uuid"))
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

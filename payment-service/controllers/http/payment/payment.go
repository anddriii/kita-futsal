package controllers

import (
	"net/http"

	errValidation "github.com/anddriii/kita-futsal/payment-service/common/error"
	"github.com/anddriii/kita-futsal/payment-service/common/response"
	"github.com/anddriii/kita-futsal/payment-service/domains/dto"
	"github.com/anddriii/kita-futsal/payment-service/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PaymentController struct {
	service service.IServiceRegistry
}

// Create implements IPaymentController.
func (p *PaymentController) Create(ctx *gin.Context) {
	var request dto.PaymentRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})
		return
	}

	result, err := p.service.GetPayment().Create(ctx, &request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: result,
		Gin:  ctx,
	})

}

// GetAllWithPagination implements IPaymentController.
func (p *PaymentController) GetAllWithPagination(ctx *gin.Context) {
	var param dto.PaymentRequestParam
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(param); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})
		return
	}

	result, err := p.service.GetPayment().GetAllWithPagination(ctx, &param)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: result,
		Gin:  ctx,
	})

}

// GetByUUID implements IPaymentController.
func (p *PaymentController) GetByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	result, err := p.service.GetPayment().GetByUUID(ctx, uuid)
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
		Data: result,
		Gin:  ctx,
	})
}

// Webhook implements IPaymentController.
func (p *PaymentController) Webhook(ctx *gin.Context) {
	var request dto.Webhook
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
	}

	err = p.service.GetPayment().WebHook(ctx, &request)
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
		Gin:  ctx,
	})
}

type IPaymentController interface {
	GetAllWithPagination(ctx *gin.Context)
	GetByUUID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Webhook(ctx *gin.Context)
}

func NewPaymentController(service service.IServiceRegistry) IPaymentController {
	return &PaymentController{
		service: service,
	}
}

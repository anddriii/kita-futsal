package controllers

import (
	"net/http"

	error2 "github.com/anddriii/kita-futsal/order-service/common/error"
	"github.com/anddriii/kita-futsal/order-service/common/response"
	"github.com/anddriii/kita-futsal/order-service/domain/dto"
	"github.com/anddriii/kita-futsal/order-service/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderController struct {
	service services.IServiceRegistry
}

// Create implements [IOrderController].
func (o *OrderController) Create(c *gin.Context) {
	var (
		request dto.OrderRequest
		ctx     = c.Request.Context()
	)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := error2.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResponse{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errorResponse,
			Gin:     c,
		})
		return
	}

	result, err := o.service.GetOrder().Create(ctx, &request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}
	response.HttpResponse(response.ParamHttpResponse{
		Code: http.StatusCreated,
		Data: result,
		Gin:  c,
	})
}

// GetAllWithPagination implements [IOrderController].
func (o *OrderController) GetAllWithPagination(c *gin.Context) {
	var params dto.OrderRequestParam
	err := c.ShouldBindQuery(&params)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	if err = validate.Struct(params); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := error2.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResponse{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errorResponse,
			Gin:     c,
		})
		return
	}

	result, err := o.service.GetOrder().GetAllWithPagination(c.Request.Context(), &params)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}
	response.HttpResponse(response.ParamHttpResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

// GetByUUID implements [IOrderController].
func (o *OrderController) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	result, err := o.service.GetOrder().GetByUUID(c.Request.Context(), uuid)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}
	response.HttpResponse(response.ParamHttpResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

// GetOrdersByUserUUID implements [IOrderController].
func (o *OrderController) GetOrdersByUserID(c *gin.Context) {

	result, err := o.service.GetOrder().GetOrderByUserID(c.Request.Context())
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}
	response.HttpResponse(response.ParamHttpResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

type IOrderController interface {
	GetAllWithPagination(*gin.Context)
	GetByUUID(*gin.Context)
	GetOrdersByUserID(*gin.Context)
	Create(*gin.Context)
}

func NewOrderController(service services.IServiceRegistry) IOrderController {
	return &OrderController{service: service}
}

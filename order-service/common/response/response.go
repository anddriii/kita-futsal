package response

import (
	"net/http"

	"github.com/anddriii/kita-futsal/order-service/constants"
	errConst "github.com/anddriii/kita-futsal/order-service/constants/error"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Token   string      `json:"token"`
}

type ParamHttpResponse struct {
	Code    int
	Err     error
	Message *string
	Gin     *gin.Context
	Data    interface{}
	Token   string
}

func HttpResponse(param ParamHttpResponse) {
	if param.Err == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Success,
			Message: http.StatusText(http.StatusOK),
			Data:    param.Data,
			Token:   param.Token,
		})
		return
	}

	message := errConst.ErrInternalServerError.Error()
	if param.Message != nil {
		message = *param.Message
	} else if param.Err != nil {
		if errConst.ErrMapping(param.Err) {
			message = param.Err.Error()
		}
	}

	param.Gin.JSON(param.Code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    param.Data,
	})
	return
}

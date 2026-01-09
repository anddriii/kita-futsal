package response

import "github.com/gin-gonic/gin"

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
		param.Gin.JSON(param.Code, Response{})
	}
}

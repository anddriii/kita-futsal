package response

import (
	"net/http"

	"github.com/anddriii/kita-futsal/user-service/constants"
	errorConstant "github.com/anddriii/kita-futsal/user-service/constants/error"
	"github.com/gin-gonic/gin"
)

// format API response
type Response struct {
	Status  string      `json:"status"`
	Message any         `json:"message"`
	Data    interface{} `json:"data"`
	Token   *string     `json:"token,omitempty"` //jika token tidak diisi maka "Token" tidak masuk ke JSON
}

type ParamHTTPResp struct {
	Code    int
	Err     error
	Message *string
	Gin     *gin.Context // Context Gin untuk mengirim response.
	Data    interface{}
	Token   *string
}

// HTTPResponse mengirim response dalam format JSON.
func HTTPResponse(param ParamHTTPResp) {
	if param.Err == nil {
		// output JSON success, no errors
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Succes,
			Message: http.StatusText(http.StatusOK),
			Data:    param.Data,
			Token:   param.Token,
		})
	}

	message := errorConstant.ErrInternalServerError.Error()
	if param.Message != nil {
		message = *param.Message
	} else if param.Err != nil {
		if errorConstant.ErrMapping(param.Err) {
			message = param.Err.Error()
		}
	}

	//ouput JSON error
	param.Gin.JSON(param.Code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    param.Data,
	})
}

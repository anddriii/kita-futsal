package controllers

import (
	"fmt"
	"net/http"

	errValidation "github.com/anddriii/kita-futsal/field-service/common/error"
	"github.com/anddriii/kita-futsal/field-service/common/response"
	"github.com/anddriii/kita-futsal/field-service/domains/dto"
	"github.com/anddriii/kita-futsal/field-service/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
)

type FieldController struct {
	service services.IServiceRegistry
}

// Create implements IFieldController and handles the creation of a Field resource.
func (f *FieldController) Create(c *gin.Context) {
	// Define a variable to hold the request payload
	var request dto.FieldRequest
	fmt.Print("sudah masuk ke controller request: ", request)

	// Bind incoming request data from multipart form into the request struct
	err := c.ShouldBindWith(&request, binding.FormMultipart)
	if err != nil {
		log.Errorf("error from controller shouldiBind", err)
		// Return a bad request response if binding fails
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	fmt.Print("sudah masuk ke controller shouldBind")

	// Initialize a validator to validate the request data
	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		log.Errorf("error from controller validate", err)
		// If validation fails, return an unprocessable entity (422) response
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errorResponse,
			Gin:     c,
		})
		return
	}

	fmt.Print("berhasil validasi: ")

	// Call the service layer to create a new Field record
	result, err := f.service.GetField().Create(c, &request)
	if err != nil {
		log.Errorf("error from controller create", err)
		// If an error occurs while creating the record, return a bad request response
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	fmt.Print("berhasil di buat", result)

	// Return a success response with the created resource
	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusCreated, // HTTP 201 Created
		Data: result,             // Response data from the service layer
		Err:  err,
		Gin:  c,
	})
}

// Delete implements IFieldController.
func (f *FieldController) Delete(c *gin.Context) {
	err := f.service.GetField().Delete(c, c.Param("uuid"))
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Gin:  c,
	})
}

// GetAllWithPagination implements IFieldController and retrieves a paginated list of Fields.
func (f *FieldController) GetAllWithPagination(c *gin.Context) {
	// Define a variable to hold query parameters from the request
	var params dto.FieldRequestParam

	// Bind query parameters from the request URL to the params struct
	err := c.ShouldBindQuery(&params)
	if err != nil {
		// Return a bad request response if query parameters binding fails
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest, // HTTP 400 Bad Request
			Err:  err,
			Gin:  c,
		})
		return
	}

	// Initialize a validator to validate query parameters
	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		// If validation fails, return an unprocessable entity (422) response
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest, // HTTP 400 Bad Request
			Err:     err,
			Message: &errMessage,
			Data:    errorResponse,
			Gin:     c,
		})
		return
	}

	// Call the service layer to retrieve paginated Field records
	result, err := f.service.GetField().GetAllWithPagination(c, &params)
	if err != nil {
		// If an error occurs while retrieving data, return a bad request response
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest, // HTTP 400 Bad Request
			Err:  err,
			Gin:  c,
		})
		return
	}

	// Return a success response with the retrieved paginated data
	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK, // HTTP 200 OK
		Data: result,        // Paginated data result from the service layer
		Gin:  c,
	})
}

// GetAllWithoutPagination implements IFieldController.
func (f *FieldController) GetAllWithoutPagination(c *gin.Context) {
	result, err := f.service.GetField().GetAllWithoutPagination(c)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

// GetByUUID implements IFieldController.
func (f *FieldController) GetByUUID(c *gin.Context) {
	result, err := f.service.GetField().GetByUUID(c, c.Param("uuid"))
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

// Update implements IFieldController.
func (f *FieldController) Update(c *gin.Context) {
	var request dto.UpdateFieldRequest
	err := c.ShouldBindWith(&request, binding.FormMultipart)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errorResponse,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetField().Update(c, c.Param("uuid"), &request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func NewFieldController(service services.IServiceRegistry) IFieldController {
	return &FieldController{service: service}
}

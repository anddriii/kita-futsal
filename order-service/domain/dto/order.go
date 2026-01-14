package dto

import (
	"time"

	"github.com/anddriii/kita-futsal/order-service/constants"
	"github.com/google/uuid"
)

type OrderRequest struct {
	FieldScheduleIDs []string `json:"fieldScheduleIDs" validate:"required,min=1,dive,required"`
}

type OrderRequestParam struct {
	Page       int     `form:"page" validate:"required"`
	Limit      int     `form:"limit" validate:"required"`
	SortColumn *string `form:"sortColumn"`
	SortOrder  *string `form:"sortOrder"`
}

type OrderResponse struct {
	UUID        uuid.UUID                   `json:"uuid"`
	Code        string                      `json:"code"`
	UserName    string                      `json:"userName"`
	Amount      float64                     `json:"amount"`
	Status      constants.OrderStatusString `json:"status"`
	PaymentLink string                      `json:"paymentLink"`
	OrderDate   time.Time                   `json:"orderDate"`
	CreatedAt   time.Time                   `json:"createdAt"`
	UpdatedAt   time.Time                   `json:"updatedAt"`
}

type OrderByUserIDResponse struct {
	Code        string                      `json:"code"`
	Amount      string                      `json:"amount"`
	Status      constants.OrderStatusString `json:"status"`
	PaymentLink string                      `json:"paymentLink"`
	OrderDate   string                      `json:"orderDate"`
	InvoiceLink *string                     `json:"invoiceLink"`
}

package dto

import "github.com/anddriii/kita-futsal/order-service/constants"

type OrderHistoryRequest struct {
	OrderID uint
	Status  constants.OrderStatusString
}

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/anddriii/kita-futsal/order-service/clients"
	clientField "github.com/anddriii/kita-futsal/order-service/clients/field"
	clientPayment "github.com/anddriii/kita-futsal/order-service/clients/payment"
	clientUser "github.com/anddriii/kita-futsal/order-service/clients/user"
	"github.com/anddriii/kita-futsal/order-service/common/util"
	"github.com/anddriii/kita-futsal/order-service/constants"
	errOrder "github.com/anddriii/kita-futsal/order-service/constants/error/order"
	"github.com/anddriii/kita-futsal/order-service/domain/dto"
	"github.com/anddriii/kita-futsal/order-service/domain/models"
	"github.com/anddriii/kita-futsal/order-service/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	repository repositories.IRepositoryRegistry
	client     clients.IClientRegistry
}

// Create implements [IOrdderService].
func (o *OrderService) Create(ctx context.Context, req *dto.OrderRequest) (*dto.OrderResponse, error) {
	var (
		order              *models.Order
		txErr, err         error
		user               = ctx.Value(constants.User).(*clientUser.UserData)
		field              *clientField.FieldData
		paymentResponse    *clientPayment.PaymentData
		orderFieldSchedule = make([]models.OrderField, 0, len(req.FieldScheduleIDs))
		totalAmount        float64
	)

	for _, fieldID := range req.FieldScheduleIDs {
		uuidParsed := uuid.MustParse(fieldID)
		field, err = o.client.GetField().GetFieldByUUID(ctx, uuidParsed)
		if err != nil {
			return nil, err
		}

		totalAmount += field.PricePerHour
		if field.Status == constants.BookesStatus.String() {
			return nil, errOrder.ErrFieldAlreadyBooked
		}
	}

	err = o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		order, txErr = o.repository.GetOrder().Create(ctx, tx, &models.Order{
			UserID: user.UUID,
			Amount: totalAmount,
			Date:   time.Now(),
			Status: constants.Pending,
			IsPaid: false,
		})
		if txErr != nil {
			return txErr
		}

		for _, fieldID := range req.FieldScheduleIDs {
			uuidParsed := uuid.MustParse(fieldID)
			orderFieldSchedule = append(orderFieldSchedule, models.OrderField{
				OrderID:         order.ID,
				FieldSCheduleID: uuidParsed,
			})
		}

		txErr = o.repository.GetOrderField().Create(ctx, tx, orderFieldSchedule)
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  constants.Pending.GetStatusString(),
			OrderID: order.ID,
		})
		if txErr != nil {
			return txErr
		}

		expiredAT := time.Now().Add(1 * time.Hour)
		description := fmt.Sprintf("Pembayaran Sewa %s", field.FieldName)
		paymentResponse, txErr = o.client.GetPayment().CreatePaymentLink(ctx, &dto.PaymentRequest{
			OrderID:     order.UUID,
			ExpiredAt:   expiredAT,
			Amount:      totalAmount,
			Description: description,
			CustomerDetail: dto.CustomerDetail{
				Name:  user.Name,
				Email: user.Email,
				Phone: user.PhoneNumber,
			},
			ItemDetails: []dto.ItemDetails{
				{
					ID:       uuid.New(),
					Name:     description,
					Amount:   totalAmount,
					Quantity: 1,
				},
			},
		})
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrder().Update(ctx, tx, &models.Order{
			PaymentID: paymentResponse.UUID,
		}, order.UUID)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		UUID:        order.UUID,
		Code:        order.Code,
		UserName:    user.Name,
		Amount:      order.Amount,
		Status:      order.Status.GetStatusString(),
		OrderDate:   order.Date,
		PaymentLink: paymentResponse.PaymentLink,
		CreatedAt:   *order.CreatedAt,
		UpdatedAt:   *order.UpdatedAt,
	}
	return &response, nil

}

// GetAllWithPagination implements [IOrdderService].
func (o *OrderService) GetAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) (*util.PaginationResult, error) {
	orders, total, err := o.repository.GetOrder().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}
	orderResults := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		user, err := o.client.GetUser().GetUserByUUID(ctx, order.UserID)
		if err != nil {
			return nil, err
		}
		orderResults = append(orderResults, dto.OrderResponse{
			UUID:      order.UUID,
			Code:      order.Code,
			UserName:  user.Name,
			Amount:    order.Amount,
			Status:    order.Status.GetStatusString(),
			OrderDate: order.Date,
			CreatedAt: *order.CreatedAt,
			UpdatedAt: *order.UpdatedAt,
		})
	}

	paginationParam := util.PaginationParam{
		Page:  param.Page,
		Limit: param.Limit,
		Count: total,
		Data:  orderResults,
	}

	response := util.GeneratePagination(paginationParam)
	return &response, nil
}

// GetByUUID implements [IOrdderService].
func (o *OrderService) GetByUUID(ctx context.Context, uuid string) (*dto.OrderResponse, error) {
	var (
		order *models.Order
		user  *clientUser.UserData
		err   error
	)

	order, err = o.repository.GetOrder().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	user, err = o.client.GetUser().GetUserByUUID(ctx, order.UserID)
	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		UUID:      order.UUID,
		Code:      order.Code,
		UserName:  user.Name,
		Amount:    order.Amount,
		Status:    order.Status.GetStatusString(),
		OrderDate: order.Date,
		CreatedAt: *order.CreatedAt,
		UpdatedAt: *order.UpdatedAt,
	}

	return &response, nil
}

// GetOrderByUserID implements [IOrdderService].
func (o *OrderService) GetOrderByUserID(ctx context.Context) ([]dto.OrderByUserIDResponse, error) {
	var (
		order []models.Order
		user  = ctx.Value(constants.User).(*clientUser.UserData)
		err   error
	)
	order, err = o.repository.GetOrder().FindByUserID(ctx, user.UUID.String())
	if err != nil {
		return nil, err
	}

	orderLists := make([]dto.OrderByUserIDResponse, 0, len(order))
	for _, item := range order {
		payment, err := o.client.GetPayment().GetPaymentByUUID(ctx, item.PaymentID)
		if err != nil {
			return nil, err
		}

		orderLists = append(orderLists, dto.OrderByUserIDResponse{
			Code:        item.Code,
			Amount:      fmt.Sprintf("%s", util.RupiahFormat(&item.Amount)),
			Status:      item.Status.GetStatusString(),
			OrderDate:   item.Date.Format(time.DateOnly),
			PaymentLink: payment.PaymentLink,
			InvoiceLink: payment.InvoiceLink,
		})
	}
	return orderLists, nil
}

func (o *OrderService) mapPaymentStatusToOrder(request *dto.PaymentData) (constants.OrderStatus, *models.Order) {
	var (
		status constants.OrderStatus
		order  *models.Order
	)

	switch request.Status {
	case constants.SettlementPaymentStatus:
		status = constants.PaymentSuccess
		order = &models.Order{
			IsPaid:    true,
			PaymentID: request.PaymentID,
			PaidAt:    request.PaidAt,
			Status:    status,
		}
	case constants.ExpirePaymentStatus:
		status = constants.Expired
		order = &models.Order{
			IsPaid:    false,
			PaymentID: request.PaymentID,
			Status:    status,
		}
	case constants.PendingPaymentStatus:
		status = constants.PendingPayment
		order = &models.Order{
			IsPaid:    false,
			PaymentID: request.PaymentID,
			Status:    status,
		}
	}
	return status, order
}

// HandlePayment implements [IOrdderService].
func (o *OrderService) HandlePayment(ctx context.Context, request *dto.PaymentData) error {
	var (
		err, txErr          error
		order               *models.Order
		orderFieldSchedules []models.OrderField
	)
	status, body := o.mapPaymentStatusToOrder(request)
	err = o.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		txErr = o.repository.GetOrder().Update(ctx, tx, body, request.OrderID)
		if txErr != nil {
			return txErr
		}

		order, txErr = o.repository.GetOrder().FindByUUID(ctx, request.OrderID.String())
		if txErr != nil {
			return txErr
		}

		txErr = o.repository.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  status.GetStatusString(),
			OrderID: order.ID,
		})

		if request.Status == constants.SettlementPaymentStatus {
			orderFieldSchedules, txErr = o.repository.GetOrderField().FindByOrderID(ctx, order.ID)
			if txErr != nil {
				return txErr
			}

			filedScheduleIDs := make([]string, 0, len(orderFieldSchedules))
			for _, item := range orderFieldSchedules {
				filedScheduleIDs = append(filedScheduleIDs, item.FieldSCheduleID.String())
			}

			txErr = o.client.GetField().UpdateStatus(&dto.UpdateFieldScheduleStatusRequest{
				FieldSCheduleIDs: filedScheduleIDs,
			})
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func NewOrderService(repository repositories.IRepositoryRegistry, client clients.IClientRegistry) IOrdderService {
	return &OrderService{repository: repository, client: client}
}

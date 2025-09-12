package routes

import (
	"github.com/anddriii/kita-futsal/payment-service/clients"
	"github.com/anddriii/kita-futsal/payment-service/constants"
	controllers "github.com/anddriii/kita-futsal/payment-service/controllers/http"
	"github.com/anddriii/kita-futsal/payment-service/middlewares"
	"github.com/gin-gonic/gin"
)

type PaymentRoute struct {
	controller controllers.IControllerRegistry
	client     clients.IClientRegistry
	group      *gin.RouterGroup
}

type IPaymentRoute interface {
	Run()
}

func NewPaymentRoute(
	group *gin.RouterGroup,
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry,
) IPaymentRoute {
	return &PaymentRoute{
		controller: controller,
		client:     client,
		group:      group,
	}
}

func (p *PaymentRoute) Run() {
	group := p.group.Group("/payment")
	group.POST("/webhook", p.controller.GetPayment().Webhook)
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
		constants.User,
	}, p.client), p.controller.GetPayment().GetAllWithPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.User,
	}, p.client), p.controller.GetPayment().GetByUUID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, p.client), p.controller.GetPayment().Create)
}

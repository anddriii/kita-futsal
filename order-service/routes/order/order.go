package routes

import (
	"github.com/anddriii/kita-futsal/order-service/clients"
	"github.com/anddriii/kita-futsal/order-service/constants"
	controllers "github.com/anddriii/kita-futsal/order-service/controllers/http"
	"github.com/anddriii/kita-futsal/order-service/middlewares"
	"github.com/gin-gonic/gin"
)

type OrderRoute struct {
	controllers.IControllerRegistry
	client clients.IClientRegistry
	group  *gin.RouterGroup
}

type IOrderRoute interface {
	Run()
}

func NewOrderRoute(group *gin.RouterGroup, client clients.IClientRegistry, controller controllers.IControllerRegistry) IOrderRoute {
	return &OrderRoute{
		IControllerRegistry: controller,
		client:              client,
		group:               group,
	}
}

func (o *OrderRoute) Run() {
	group := o.group.Group("/orders")
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, o.client), o.GetOrder().GetAllWithPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, o.client), o.GetOrder().GetByUUID)
	group.GET("/user", middlewares.CheckRole([]string{
		constants.Customer,
	}, o.client), o.GetOrder().GetOrdersByUserID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Customer,
	}, o.client), o.GetOrder().Create)
}

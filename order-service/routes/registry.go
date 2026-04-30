package routes

import (
	"github.com/anddriii/kita-futsal/order-service/clients"
	controllers "github.com/anddriii/kita-futsal/order-service/controllers/http"
	routes "github.com/anddriii/kita-futsal/order-service/routes/order"
	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	client     clients.IClientRegistry
	group      *gin.RouterGroup
}

type IRouteRegistry interface {
	Serve()
}

func NewRouteRegistry(
	group *gin.RouterGroup,
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry,
) IRouteRegistry {
	return &Registry{
		controller: controller,
		client:     client,
		group:      group,
	}
}

func (r *Registry) Serve() {
	r.orderRoute().Run()
}

func (r *Registry) orderRoute() routes.IOrderRoute {
	return routes.NewOrderRoute(r.group, r.client, r.controller)
}

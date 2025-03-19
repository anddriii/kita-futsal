package routes

import (
	"github.com/anddriii/kita-futsal/field-service/clients"
	"github.com/anddriii/kita-futsal/field-service/controllers"
	"github.com/gin-gonic/gin"

	fieldRoute "github.com/anddriii/kita-futsal/field-service/routes/field"
	fieldScheduleRoute "github.com/anddriii/kita-futsal/field-service/routes/field_schedule"
	timeRoute "github.com/anddriii/kita-futsal/field-service/routes/time"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IRegistry interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IRegistry {
	return &Registry{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (r *Registry) fieldRoute() fieldRoute.IFieldRoute {
	return fieldRoute.NewFieldRoute(r.controller, r.group, r.client)
}

func (r *Registry) fieldScheduleRoute() fieldScheduleRoute.IFieldScheduleRoute {
	return fieldScheduleRoute.NewFieldScheduleRoute(r.controller, r.group, r.client)
}

func (r *Registry) timeRoute() timeRoute.ITimeRoute {
	return timeRoute.NewRouteTime(r.controller, r.group, r.client)
}

func (r *Registry) Serve() {
	r.fieldRoute().Run()
	r.fieldScheduleRoute().Run()
	r.timeRoute().Run()
}

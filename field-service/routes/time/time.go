package routes

import (
	"github.com/anddriii/kita-futsal/field-service/clients"
	"github.com/anddriii/kita-futsal/field-service/constants"
	"github.com/anddriii/kita-futsal/field-service/controllers"
	"github.com/anddriii/kita-futsal/field-service/middlewares"
	"github.com/gin-gonic/gin"
)

type TimeRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

// Run implements ITimeRoute.
func (t *TimeRoute) Run() {
	group := t.group.Group("/time")
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
	}, t.client),
		t.controller.GetTime().GetAll)

	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, t.client),
		t.controller.GetTime().GetByUUID)

	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, t.client),
		t.controller.GetTime().Create)
}

type ITimeRoute interface {
	Run()
}

func NewRouteTime(controler controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) ITimeRoute {
	return &TimeRoute{
		controller: controler,
		group:      group,
		client:     client,
	}
}

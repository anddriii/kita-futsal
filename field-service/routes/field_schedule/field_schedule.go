package routes

import (
	"github.com/anddriii/kita-futsal/field-service/clients"
	"github.com/anddriii/kita-futsal/field-service/constants"
	"github.com/anddriii/kita-futsal/field-service/controllers"
	"github.com/anddriii/kita-futsal/field-service/middlewares"
	"github.com/gin-gonic/gin"
)

type FieldScheduleRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IFieldScheduleRoute interface {
	Run()
}

func NewFieldScheduleRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IFieldScheduleRoute {
	return &FieldScheduleRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

// Run implements IFieldScheduleRoute and sets up all field schedule-related routes.
func (f FieldScheduleRoute) Run() {
	// Create a sub-route group for field schedule management
	group := f.group.Group("/field/schedule")

	// Get a list of schedules by field ID and date (no authentication token required)
	group.GET("/lists/:uuid", middlewares.AuthenticateWithoutToken(),
		f.controller.GetFieldSchedule().GetAllByFieldIdAndDate)

	// Update schedule status (no authentication token required)
	group.PATCH("/status", middlewares.AuthenticateWithoutToken(),
		f.controller.GetFieldSchedule().UpdateStatus)

	// Apply authentication middleware for routes below
	group.Use(middlewares.Authenticate())

	// Get paginated schedule list (accessible by Admin & User roles)
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Admin,
		constants.User,
	}, f.client), f.controller.GetFieldSchedule().GetAllWithPagination)

	// Get schedule details by UUID (accessible by Admin & User roles)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.User,
	}, f.client), f.controller.GetFieldSchedule().GetByUUID)

	// Create a new schedule (only Admin can access)
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetFieldSchedule().Create)

	// Generate schedule for one month (only Admin can access)
	group.POST("/one-month", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetFieldSchedule().GenerateScheduleForOneMonth)

	// Update an existing schedule by UUID (only Admin can access)
	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetFieldSchedule().Update)

	// Delete a schedule by UUID (only Admin can access)
	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetFieldSchedule().Delete)
}

package routes

import (
	"github.com/anddriii/kita-futsal/field-service/clients"
	"github.com/anddriii/kita-futsal/field-service/constants"
	"github.com/anddriii/kita-futsal/field-service/controllers"
	"github.com/anddriii/kita-futsal/field-service/middlewares"
	"github.com/gin-gonic/gin"
)

type FieldRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IFieldRoute interface {
	Run()
}

func NewFieldRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IFieldRoute {
	return &FieldRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (f *FieldRoute) Run() {
	group := f.group.Group("/field")

	//endpoint without login
	group.GET("", middlewares.AuthenticateWithoutToken(), f.controller.GetField().GetAllWithoutPagination)
	group.GET("/:uuid", middlewares.AuthenticateWithoutToken(), f.controller.GetField().GetByUUID)
	group.GET("/nearby", middlewares.AuthenticateWithoutToken(), f.controller.GetField().GetNearbyFields)

	//Middleware autentikasi diterapkan ke seluruh route berikutnya
	group.Use(middlewares.Authenticate())

	//endpoint must login

	// Mengambil semua field dengan pagination, hanya bisa diakses oleh Admin & User
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Admin,
		constants.User,
	}, f.client),
		f.controller.GetField().GetAllWithPagination)

	// Membuat field baru, hanya bisa diakses oleh Admin
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client), f.controller.GetField().Create)

	// Memperbarui field berdasarkan UUID, hanya bisa diakses oleh Admin
	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client),
		f.controller.GetField().Update)

	// menghapus field beradasarkan UUID, hanya bisa diakses oleh admin
	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, f.client),
		f.controller.GetField().Delete)
}

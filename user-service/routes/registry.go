package routes

import (
	"github.com/anddriii/kita-futsal/user-service/controllers"
	routes "github.com/anddriii/kita-futsal/user-service/routes/user"
	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	userGroup  *gin.RouterGroup
}

type IRouteRegister interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup) IRouteRegister {
	return &Registry{controller: controller, userGroup: group}
}

// Serve implements IRouteRegister.
func (r *Registry) Serve() {
	r.userRoute().Run()
}

func (r *Registry) userRoute() routes.IUserRoute {
	return routes.NewUserRoute(r.controller, r.userGroup)
}

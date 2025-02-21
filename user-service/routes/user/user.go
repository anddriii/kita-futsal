package routes

import (
	"github.com/anddriii/kita-futsal/user-service/controllers"
	"github.com/anddriii/kita-futsal/user-service/middlewares"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controller controllers.IControllerRegistry
	userGroup  *gin.RouterGroup
}

type IUserRoute interface {
	Run()
}

func NewUserRoute(contoller controllers.IControllerRegistry, group *gin.RouterGroup) IUserRoute {
	return &UserRoute{controller: contoller, userGroup: group}
}

// Run implements IUserRoute.
func (u *UserRoute) Run() {
	group := u.userGroup.Group("/auth")
	group.GET("/user", middlewares.Authenticate(), u.controller.GetUserController().GetUserLogin)
	group.GET("/:uuid", middlewares.Authenticate(), u.controller.GetUserController().GetUserUUID)
	group.POST("/login", u.controller.GetUserController().Login)
	group.POST("/register", u.controller.GetUserController().Register)
	group.PUT("/:uuid", middlewares.Authenticate(), u.controller.GetUserController().Update)
}

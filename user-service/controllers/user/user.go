package controllers

import "github.com/gin-gonic/gin"

type IUserController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetUserLogin(ctx *gin.Context)
	GetUserUUID(ctx *gin.Context)
}

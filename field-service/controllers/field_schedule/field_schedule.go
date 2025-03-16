package controllers

import "github.com/gin-gonic/gin"

type IFieldScheduleController interface {
	GetAllWithPagination(ctx *gin.Context)
	GetAllByFieldIdAndDate(ctx *gin.Context)
	GetByUUID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	UpdateStatus(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GenerateScheduleForOneMonth(ctx *gin.Context)
}

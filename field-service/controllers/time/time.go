package controllers

import "github.com/gin-gonic/gin"

type ITimeController interface {
	GetAll(c *gin.Context)
	GetByUUID(c *gin.Context)
	Create(c *gin.Context)
}

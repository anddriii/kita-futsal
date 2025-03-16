package controllers

import (
	fieldController "github.com/anddriii/kita-futsal/field-service/controllers/field"
	fieldScheduleController "github.com/anddriii/kita-futsal/field-service/controllers/field_schedule"
	timeController "github.com/anddriii/kita-futsal/field-service/controllers/time"
	"github.com/anddriii/kita-futsal/field-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

// GetField implements IControllerRegistry.
func (r *Registry) GetField() fieldController.IFieldController {
	return fieldController.NewFieldController(r.service)
}

// GetFieldSchedule implements IControllerRegistry.
func (r *Registry) GetFieldSchedule() fieldScheduleController.IFieldScheduleController {
	return fieldScheduleController.NewFieldScheduleController(r.service)
}

// GetTime implements IControllerRegistry.
func (r *Registry) GetTime() timeController.ITimeController {
	return timeController.NewTimeController(r.service)
}

type IControllerRegistry interface {
	GetField() fieldController.IFieldController
	GetFieldSchedule() fieldScheduleController.IFieldScheduleController
	GetTime() timeController.ITimeController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

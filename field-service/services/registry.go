package services

import (
	"github.com/anddriii/kita-futsal/field-service/common/gcs"
	"github.com/anddriii/kita-futsal/field-service/repositories"
	fieldService "github.com/anddriii/kita-futsal/field-service/services/field"
	fieldScheduleService "github.com/anddriii/kita-futsal/field-service/services/field_schedule"
	timeService "github.com/anddriii/kita-futsal/field-service/services/time"
)

type Registry struct {
	repository repositories.IRepoRegistry
	gcs        gcs.IGCSClient
}

// GetField implements IServiceRegistry.
func (r *Registry) GetField() fieldService.IFieldService {
	return fieldService.NewFieldService(r.repository, r.gcs)
}

// GetFieldSchedule implements IServiceRegistry.
func (r *Registry) GetFieldSchedule() fieldScheduleService.IFieldScheduleService {
	return fieldScheduleService.NewFieldScheduleService(r.repository)
}

// GetTime implements IServiceRegistry.
func (r *Registry) GetTime() timeService.ITimeService {
	return timeService.NewTimeService(r.repository)
}

type IServiceRegistry interface {
	GetField() fieldService.IFieldService
	GetFieldSchedule() fieldScheduleService.IFieldScheduleService
	GetTime() timeService.ITimeService
}

func NewServiceRegistry(repository repositories.IRepoRegistry, gcs gcs.IGCSClient) IServiceRegistry {
	return &Registry{
		repository: repository,
		gcs:        gcs,
	}
}

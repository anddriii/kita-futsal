package repositories

import (
	fieldRepo "github.com/anddriii/kita-futsal/field-service/repositories/field"
	fieldSchedu "github.com/anddriii/kita-futsal/field-service/repositories/field_schedule"
	fieldTime "github.com/anddriii/kita-futsal/field-service/repositories/time"
	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

// GetField implements IRepoRegistry.
func (r *Registry) GetField() fieldRepo.IFieldRepository {
	return fieldRepo.NewFieldRepository(r.db)
}

// GetFieldSchedule implements IRepoRegistry.
func (r *Registry) GetFieldSchedule() fieldSchedu.IFieldScheduleRepository {
	return fieldSchedu.NewFieldScheduleRepository(r.db)
}

// GetTime implements IRepoRegistry.
func (r *Registry) GetTime() fieldTime.ITimeRepository {
	return fieldTime.NewTimeRepository(r.db)
}

type IRepoRegistry interface {
	GetField() fieldRepo.IFieldRepository
	GetFieldSchedule() fieldSchedu.IFieldScheduleRepository
	GetTime() fieldTime.ITimeRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepoRegistry {
	return &Registry{db: db}
}

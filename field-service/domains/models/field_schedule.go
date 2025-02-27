package models

import (
	"time"

	"github.com/anddriii/kita-futsal/field-service/constants"
	"github.com/google/uuid"
)

type FieldSchedule struct {
	ID        uint                          `gorm:"primaryKey;autoIncrement;not null"`
	UUID      uuid.UUID                     `gorm:"type:uuid;not null"`
	FieldId   uint                          `gorm:"type:int;not null"`
	TimeId    uint                          `gorm:"type:int;not null"`
	Date      time.Time                     `gorm:"type:date;not null"`
	Status    constants.FieldScheduleStatus `gorm:"type:int; not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
	Field     Field `gorm:"foreignKey:field_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Time      Time  `gorm:"foreignKey:time_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

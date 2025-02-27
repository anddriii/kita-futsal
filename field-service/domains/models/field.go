package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Field struct {
	ID            uint           `gorm:"primaryKey;autoIncrement;not null"`
	UUID          uuid.UUID      `gorm:"type:uuid;not null"`
	Code          string         `gorm:"type:varchar(15);not null"`
	Name          string         `gorm:"type:varchar(25);not null"`
	Description   string         `gorm:"type:text;not null"`
	Image         pq.StringArray `gorm:"type:text[];not null"`
	PricePerHour  int            `gorm:"type:int;not null"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	FieldSchedule []FieldSchedule `gorm:"foreignKey:field_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

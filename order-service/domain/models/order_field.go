package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderField struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	OrderID         uint      `gorm:"column:order_id;type:bigint;not null"`
	FieldScheduleID uuid.UUID `gorm:"column:field_schedule_id;type:uuid;not null"`
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

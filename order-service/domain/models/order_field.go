package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderField struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	OrderField      uint      `gorm:"type:bigint;not null"`
	FieldSCheduleID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

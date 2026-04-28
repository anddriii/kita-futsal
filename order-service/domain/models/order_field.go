package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderField struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	OrderID         uint      `gorm:"type:bigint;not null"`
<<<<<<< HEAD
	FieldSCheduleID uuid.UUID `gorm:"type:uuid;not null"`
=======
	FieldScheduleID uuid.UUID `gorm:"type:uuid;not null"`
>>>>>>> ad31e98 (feat(order-service): update KafkaMessage struct fields and enhance OrderRequest validation; add middleware for panic handling and rate limiting)
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

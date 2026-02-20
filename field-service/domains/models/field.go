package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Field struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;not null"`
	UUID        uuid.UUID `gorm:"type:uuid;not null"`
	Code        string    `gorm:"type:varchar(15);not null"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text;not null"`
	// Image       datatypes.JSON `gorm:"type:json;not null"`
	Image         pq.StringArray `gorm:"type:text[];not null"`
	Latitude      float64        `gorm:"type:decimal(10,8);not null"`
	Lonitude      float64        `gorm:"type:decimal(11,8);not null"`
	PricePerHour  int            `gorm:"type:int;not null"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	FieldSchedule []FieldSchedule `gorm:"foreignKey:field_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

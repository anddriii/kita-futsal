package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;not null"`
	UUID      uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Username  string    `gorm:"type:varchar(20);not null"`
	Password  string    `gorm:"type:varchar(255):not null"`
	Email     string    `gorm:"type:varchar(100);not null"`
	RoleId    uint      `gorm:"type:uint;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Role      Role `gorm:"foreignKey:role_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

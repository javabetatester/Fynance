package investment

import (
	"time"

	"github.com/google/uuid"
)

type Investment struct {
	Id              int       `gorm:"primaryKey;autoIncrement"`
	UserId          uuid.UUID `gorm:"type:uuid;index;not null"`
	Type            Types     `gorm:"not null"`
	Name            string    `gorm:"size:100;not null"`
	TotalAmount     float64   `gorm:"not null"`
	CurrentAmount   float64   `gorm:"not null"`
	ReturnRate      float64   `gorm:"not null"`
	ApplicationDate time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

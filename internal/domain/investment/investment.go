package investment

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Investment struct {
	Id              ulid.ULID `gorm:"type:varchar(26);primaryKey;serializer:json"`
	UserId          ulid.ULID `gorm:"type:varchar(26);index;not null;serializer:json"`
	Type            Types     `gorm:"not null"`
	Name            string    `gorm:"size:100;not null"`
	TotalAmount     float64   `gorm:"not null"`
	CurrentAmount   float64   `gorm:"not null"`
	ReturnRate      float64   `gorm:"not null"`
	ApplicationDate time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

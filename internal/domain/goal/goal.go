package goal

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	Id            uuid.UUID `gorm:"primaryKey"`
	UserId        uuid.UUID `gorm:"type:uuid;index;not null"`
	Name          string    `gorm:"not null"`
	TargetAmount  float64   `gorm:"not null"`
	CurrentAmount float64   `gorm:"not null"`
	StartedAt     time.Time
	EndedAt       time.Time
	Status        GoalStatus `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

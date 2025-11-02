package goal

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Goal struct {
	Id            ulid.ULID `gorm:"type:varchar(26);primaryKey"`
	UserId        ulid.ULID `gorm:"type:varchar(26);index;not null"`
	Name          string    `gorm:"not null"`
	TargetAmount  float64   `gorm:"not null"`
	CurrentAmount float64   `gorm:"not null"`
	StartedAt     time.Time
	EndedAt       time.Time
	Status        GoalStatus `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

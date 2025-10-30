package goal

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	Id            int
	UserId        uuid.UUID
	Name          string
	TargetAmount  float64
	CurrentAmount float64
	StartedAt     time.Time
	EndedAt       time.Time
	Status        GoalStatus
	CreatedAt     time.Time
}

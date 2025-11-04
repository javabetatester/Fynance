package goal

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type GoalCreateRequest struct {
	UserId     ulid.ULID  `json:"user_id"`
	Name       string     `json:"name"`
	Target     float64    `json:"target"`
	EndedAt    *time.Time `json:"end_at"` // Nullable
}

type GoalUpdateRequest struct {
	Id         ulid.ULID  `json:"id"`
	UserId     ulid.ULID  `json:"user_id"`
	Name       string     `json:"name"`
	Target     float64    `json:"target"`
	EndedAt    *time.Time `json:"end_at"` // Nullable
}

type GoalDashboardResponse struct {
	TargetAmount  float64    `json:"target_amount"`
	CurrentAmount float64    `json:"current_amount"`
	EndedAt       *time.Time `json:"ended_at"` // Nullable
	Status        GoalStatus `json:"status"`
}

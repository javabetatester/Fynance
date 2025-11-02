package goal

import "time"

type GoalCreateRequest struct {
	Name       string    `json:"name"`
	Target     float64   `json:"target"`
	CategoryId string    `json:"category_id"`
	EndedAt    time.Time `json:"end_at"`
}

type GoalDashboardResponse struct {
	TargetAmount  float64 `json:"progress_initial_amount"`
	CurrentAmount float64 `json:"progress_end_amount"`
	EndedAt       time.Time
	Status        GoalStatus
}

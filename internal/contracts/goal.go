package contracts

import (
	"time"

	domainGoal "Fynance/internal/domain/goal"
)

type GoalCreateRequest struct {
	Name   string     `json:"name" binding:"required"`
	Target float64    `json:"target" binding:"required,gt=0"`
	EndAt  *time.Time `json:"end_at"`
}

type GoalUpdateRequest struct {
	Name   string     `json:"name" binding:"required"`
	Target float64    `json:"target" binding:"required,gt=0"`
	EndAt  *time.Time `json:"end_at"`
}

type GoalResponse struct {
	Goal *domainGoal.Goal `json:"goal"`
}

type GoalListResponse struct {
	Goals []*domainGoal.Goal `json:"goals"`
	Total int                `json:"total"`
}

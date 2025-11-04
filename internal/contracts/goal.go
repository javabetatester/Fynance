package contracts

import "time"

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

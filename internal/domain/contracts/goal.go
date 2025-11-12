package domaincontracts

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type GoalCreateRequest struct {
	UserId  ulid.ULID  `json:"user_id"`
	Name    string     `json:"name"`
	Target  float64    `json:"target"`
	EndedAt *time.Time `json:"end_at"`
}

type GoalUpdateRequest struct {
	Id      ulid.ULID  `json:"id"`
	UserId  ulid.ULID  `json:"user_id"`
	Name    string     `json:"name"`
	Target  float64    `json:"target"`
	EndedAt *time.Time `json:"end_at"`
}

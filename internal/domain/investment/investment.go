package investment

import (
	"time"

	"github.com/google/uuid"
)

type Investment struct {
	Id              int
	UserId          uuid.UUID
	Type            Types
	Name            string
	TotalAmount     float64
	CurrentAmount   float64
	ReturnRate      float64
	ApplicationDate time.Time
	UpdatedAt       time.Time
}

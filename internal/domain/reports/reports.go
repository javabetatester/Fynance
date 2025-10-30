package reports

import (
	"time"

	"github.com/google/uuid"
)

type Reports struct {
	Id          int
	UserId      uuid.UUID
	StartedAt   time.Time
	EndedAt     time.Time
	GeneratedAt time.Time
}

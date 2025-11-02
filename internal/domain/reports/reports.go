package reports

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Reports struct {
	Id          int
	UserId      ulid.ULID
	StartedAt   time.Time
	EndedAt     time.Time
	GeneratedAt time.Time
}

package transaction

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	UserId      uuid.UUID
	Id          uuid.UUID
	Type        Types
	CategoryId  int
	Amount      float64
	Description string
	Date        time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Category struct {
	UserId uuid.UUID
	Id     int
	Name   string
	Icon   string
}

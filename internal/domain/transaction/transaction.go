package transaction

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	UserId      uuid.UUID `gorm:"type:uuid;index;not null"`
	Id          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Type        Types     `gorm:"not null"`
	CategoryId  uuid.UUID `gorm:"index"`
	Amount      float64   `gorm:"not null"`
	Description string    `gorm:"size:255"`
	Date        time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

type Category struct {
	UserId uuid.UUID `gorm:"type:uuid;index;not null"`
	Id     uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name   string    `gorm:"size:100;not null"`
	Icon   string    `gorm:"size:50"`
}

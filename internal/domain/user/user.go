package user

import (
	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/investment"
	"Fynance/internal/domain/transaction"
	"time"

	"github.com/google/uuid"
)

	type User struct {
		Id           uuid.UUID                `gorm:"primaryKey;type:uuid"`
		Name         string                   `gorm:"size:100;not null"`
		Email        string                   `gorm:"size:100;uniqueIndex;not null"`
		Password     string                   `gorm:"size:255;not null"`
		CreatedAt    time.Time                `gorm:"not null"`
		UpdatedAt    time.Time                `gorm:"not null"`
		Transactions []transaction.Transaction `gorm:"foreignKey:UserId"`
		Goals        []goal.Goal              `gorm:"foreignKey:UserId"`
		Investments  []investment.Investment   `gorm:"foreignKey:UserId"`
	}

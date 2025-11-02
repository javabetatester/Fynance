package user

import (
	"time"
)

type User struct {
	Id        string    `gorm:"type:varchar(26);primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Plan      Plan      `gorm:"type:varchar(10);default:'FREE'"`
	PlanSince time.Time `gorm:"type:timestamp;default:now()"`
}

type Plan string

const (
	PlanFree  Plan = "FREE"
	PlanBasic Plan = "BASIC"
	PlanPro   Plan = "PRO"
)

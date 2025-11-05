package user

import (
	"time"
)

type User struct {
	Id        string    `gorm:"type:varchar(26);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex:idx_users_email;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null" json:"updated_at"`
	Plan      Plan      `gorm:"type:varchar(10);default:'FREE';index:idx_users_plan" json:"plan"`
	PlanSince time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"plan_since"`
}

func (User) TableName() string {
	return "users"
}

type Plan string

const (
	PlanFree  Plan = "FREE"
	PlanBasic Plan = "BASIC"
	PlanPro   Plan = "PRO"
)

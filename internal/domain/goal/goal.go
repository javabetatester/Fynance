package goal

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Goal struct {
	Id            ulid.ULID  `gorm:"type:varchar(26);primaryKey" json:"id"`
	UserId        ulid.ULID  `gorm:"type:varchar(26);index:idx_goals_user_id;not null" json:"user_id"`
	Name          string     `gorm:"type:varchar(100);not null;index:idx_goals_user_name,unique" json:"name"`
	TargetAmount  float64    `gorm:"type:decimal(15,2);not null" json:"target_amount"`
	CurrentAmount float64    `gorm:"type:decimal(15,2);not null;default:0" json:"current_amount"`
	StartedAt     time.Time  `gorm:"type:timestamp" json:"started_at"`
	EndedAt       *time.Time `gorm:"type:timestamp" json:"ended_at"` // Nullable
	Status        GoalStatus `gorm:"type:varchar(20);default:'ACTIVE';index:idx_goals_status" json:"status"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;not null" json:"updated_at"`
}

func (Goal) TableName() string {
	return "goals"
}

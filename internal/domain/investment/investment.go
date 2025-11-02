package investment

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Investment struct {
	Id              ulid.ULID `gorm:"type:varchar(26);primaryKey" json:"id"`
	UserId          ulid.ULID `gorm:"type:varchar(26);index;not null" json:"user_id"`
	Type            Types     `gorm:"type:varchar(20);not null" json:"type"`
	Name            string    `gorm:"size:100;not null" json:"name"`
	CurrentBalance  float64   `gorm:"not null;default:0" json:"current_balance"`
	ReturnRate      float64   `gorm:"default:0" json:"return_rate"`
	ApplicationDate time.Time `gorm:"not null" json:"application_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
